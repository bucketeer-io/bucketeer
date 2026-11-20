// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authz

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// Mode is what the interceptor does with the decision for a service.
type Mode int

const (
	// Shadow runs the handler as before and logs where the decision differs from its own check.
	Shadow Mode = iota
	// Enforce rejects the request on a denial and hands the editor to the handler through WithEditor.
	Enforce
)

// Interceptor authorizes every RPC of a server with Table 2: it enforces the
// decision for the services listed in enforced and shadows the others.
type Interceptor struct {
	authorizer *Authorizer
	enforced   map[string]bool
	logger     *zap.Logger
}

func NewInterceptor(authorizer *Authorizer, enforced []string, logger *zap.Logger) *Interceptor {
	set := make(map[string]bool, len(enforced))
	for _, service := range enforced {
		set[service] = true
	}
	return &Interceptor{
		authorizer: authorizer,
		enforced:   set,
		logger:     logger.Named("authz"),
	}
}

// mode is Enforce for the methods of an enforced service.
func (i *Interceptor) mode(method string) Mode {
	if i.enforced[serviceOf(method)] {
		return Enforce
	}
	return Shadow
}

// serviceOf returns the proto service of a full method name "/package.Service/Method".
func serviceOf(method string) string {
	service, _, _ := strings.Cut(strings.TrimPrefix(method, "/"), "/")
	return service
}

// CallerFromContext builds the caller from the verified access token; false for an anonymous request.
func CallerFromContext(ctx context.Context) (Caller, bool) {
	t, ok := rpc.GetAccessToken(ctx)
	if !ok || t == nil {
		return Caller{}, false
	}
	caller := Caller{Email: t.Email}
	if t.IsServiceToken {
		// Service tokens carry no name yet (RFC 0050 step 2).
		caller.Service = LegacyService
	}
	return caller.Relayed(ctx), true
}

// Unary returns the interceptor for unary RPCs.
func (i *Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if skipped(info.FullMethod) {
			return handler(ctx, req)
		}
		caller, _ := CallerFromContext(ctx)
		editor, authzErr := i.authorize(ctx, caller, info.FullMethod, req)
		if i.mode(info.FullMethod) == Enforce {
			if authzErr != nil {
				return nil, authzErr
			}
			return handler(WithEditor(ctx, editor), req)
		}
		resp, err := handler(ctx, req)
		i.compare(caller, info.FullMethod, authzErr, err)
		return resp, err
	}
}

// Stream returns the interceptor for streaming RPCs; the decision is made on the first message received.
func (i *Interceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if skipped(info.FullMethod) {
			return handler(srv, ss)
		}
		stream := &authorizedStream{ServerStream: ss, interceptor: i, method: info.FullMethod}
		err := handler(srv, stream)
		if stream.denied != nil {
			return stream.denied
		}
		if i.mode(info.FullMethod) == Shadow && stream.authorized {
			i.compare(stream.caller, info.FullMethod, stream.authzErr, err)
		}
		return err
	}
}

type authorizedStream struct {
	grpc.ServerStream
	interceptor *Interceptor
	method      string
	ctx         context.Context
	caller      Caller
	authorized  bool
	authzErr    error
	denied      error
}

func (s *authorizedStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return s.ServerStream.Context()
}

func (s *authorizedStream) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if s.authorized {
		return nil
	}
	s.authorized = true
	ctx := s.ServerStream.Context()
	s.caller, _ = CallerFromContext(ctx)
	editor, err := s.interceptor.authorize(ctx, s.caller, s.method, m)
	s.authzErr = err
	if s.interceptor.mode(s.method) != Enforce {
		return nil
	}
	if err != nil {
		s.denied = err
		return err
	}
	s.ctx = WithEditor(ctx, editor)
	return nil
}

func (i *Interceptor) authorize(
	ctx context.Context,
	caller Caller,
	method string,
	req interface{},
) (*eventproto.Editor, error) {
	msg, ok := req.(proto.Message)
	if !ok {
		return nil, ErrInternal
	}
	return i.authorizer.Authorize(ctx, caller, method, msg)
}

// compare logs a decision that disagrees with the handler's success, Unauthenticated or PermissionDenied.
func (i *Interceptor) compare(caller Caller, method string, authzErr, handlerErr error) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("email", caller.Email),
		zap.String("service", caller.Service),
		zap.Bool("publicAPI", caller.PublicAPI != nil),
	}
	if errors.Is(authzErr, ErrNoRule) {
		i.logger.Error("No rule for method", fields...)
		return
	}
	authzDenied := authzErr != nil
	if authzDenied {
		fields = append(fields, zap.Error(authzErr))
		if status.Code(authzErr) == codes.Internal {
			i.logger.Error("Rule failed", fields...)
			return
		}
	}
	switch status.Code(handlerErr) {
	case codes.OK:
		if authzDenied {
			i.logger.Warn("Would deny", fields...)
		}
	case codes.Unauthenticated, codes.PermissionDenied:
		if !authzDenied {
			i.logger.Warn("Would allow", append(fields, zap.Error(handlerErr))...)
		}
	}
}

// skipped is true for the methods gRPC itself serves, such as the health check.
func skipped(method string) bool {
	return strings.HasPrefix(method, "/grpc.")
}
