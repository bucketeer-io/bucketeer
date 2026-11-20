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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const healthCheck = "/grpc.health.v1.Health/Check"

var (
	errDenied   = status.Error(codes.PermissionDenied, "handler: denied")
	errNotFound = status.Error(codes.NotFound, "handler: not found")
)

// userCtx is a context after the auth interceptor verified a user's token.
func userCtx(email string) context.Context {
	return context.WithValue(context.Background(), rpc.AccessTokenKey, &token.AccessToken{Email: email})
}

// serviceCtx is a context after the auth interceptor verified a service
// token, with the metadata the public API gateway relays when maintainer is
// set.
func serviceCtx(maintainer string) context.Context {
	ctx := context.WithValue(context.Background(), rpc.AccessTokenKey, &token.AccessToken{
		Email:          "service@example.com",
		IsServiceToken: true,
		IsSystemAdmin:  true,
	})
	if maintainer == "" {
		return ctx
	}
	return metadata.NewIncomingContext(ctx, metadata.Pairs(
		APIKeyTokenMDKey, rawAPIKey,
		APIKeyMaintainerMDKey, maintainer,
		APIKeyNameMDKey, "key",
	))
}

type loggedEntry struct {
	level   zapcore.Level
	message string
}

func newTestInterceptor(t *testing.T, mode Mode) (*Interceptor, *observer.ObservedLogs) {
	t.Helper()
	core, logs := observer.New(zapcore.DebugLevel)
	accounts := &fakeAccounts{accounts: testAccounts()}
	authorizer := New(testRules(guardFunc(func(context.Context, *Access) error { return nil })), accounts)
	var enforced []string
	if mode == Enforce {
		enforced = []string{FeatureService}
	}
	return NewInterceptor(authorizer, enforced, zap.New(core)), logs
}

func entries(logs *observer.ObservedLogs) []loggedEntry {
	var out []loggedEntry
	for _, e := range logs.All() {
		out = append(out, loggedEntry{level: e.Level, message: e.Message})
	}
	return out
}

func TestUnaryInterceptor(t *testing.T) {
	t.Parallel()
	req := &featureproto.GetFeatureRequest{EnvironmentId: "env-1"}
	patterns := []struct {
		desc            string
		mode            Mode
		method          string
		ctx             context.Context
		req             interface{}
		handlerErr      error
		expectedHandled bool
		expectedEditor  *eventproto.Editor
		expectedCode    codes.Code
		expectedLogs    []loggedEntry
	}{
		{
			desc: "shadow: both allow", mode: Shadow, method: getFeature, ctx: userCtx("viewer@example.com"), req: req,
			expectedHandled: true,
		},
		{
			desc: "shadow: would deny", mode: Shadow, method: getFeature, ctx: userCtx("member@example.com"), req: req,
			expectedHandled: true,
			expectedLogs:    []loggedEntry{{level: zapcore.WarnLevel, message: "Would deny"}},
		},
		{
			desc: "shadow: would allow", mode: Shadow, method: getFeature, ctx: userCtx("viewer@example.com"), req: req,
			handlerErr: errDenied, expectedHandled: true, expectedCode: codes.PermissionDenied,
			expectedLogs: []loggedEntry{{level: zapcore.WarnLevel, message: "Would allow"}},
		},
		{
			desc: "shadow: both deny", mode: Shadow, method: getFeature, ctx: userCtx("member@example.com"), req: req,
			handlerErr: errDenied, expectedHandled: true, expectedCode: codes.PermissionDenied,
		},
		{
			desc: "shadow: handler error says nothing", mode: Shadow, method: getFeature,
			ctx: userCtx("member@example.com"), req: req,
			handlerErr: errNotFound, expectedHandled: true, expectedCode: codes.NotFound,
		},
		{
			desc: "shadow: no rule", mode: Shadow, method: "/bucketeer.feature.FeatureService/Unknown",
			ctx: userCtx("viewer@example.com"), req: req, expectedHandled: true,
			expectedLogs: []loggedEntry{{level: zapcore.ErrorLevel, message: "No rule for method"}},
		},
		{
			desc: "shadow: rule failed", mode: Shadow, method: invalidRule, ctx: userCtx("viewer@example.com"), req: req,
			expectedHandled: true,
			expectedLogs:    []loggedEntry{{level: zapcore.ErrorLevel, message: "Rule failed"}},
		},
		{
			desc: "shadow: public", mode: Shadow, method: ping, ctx: context.Background(), req: req,
			expectedHandled: true,
		},
		{
			desc: "shadow: health check is skipped", mode: Shadow, method: healthCheck, ctx: context.Background(), req: req,
			expectedHandled: true,
		},
		{
			desc: "enforce: allow", mode: Enforce, method: getFeature, ctx: userCtx("viewer@example.com"), req: req,
			expectedHandled: true, expectedEditor: memberEditor("viewer@example.com"),
		},
		{
			desc: "enforce: deny", mode: Enforce, method: getFeature, ctx: userCtx("member@example.com"), req: req,
			expectedCode: codes.PermissionDenied,
		},
		{
			desc: "enforce: no token", mode: Enforce, method: getFeature, ctx: context.Background(), req: req,
			expectedCode: codes.Unauthenticated,
		},
		{
			desc: "enforce: service", mode: Enforce, method: getFeature, ctx: serviceCtx(""), req: req,
			expectedHandled: true, expectedEditor: &eventproto.Editor{Email: "service@example.com", IsAdmin: true},
		},
		{
			desc: "enforce: relayed as the maintainer", mode: Enforce, method: getFeature,
			ctx: serviceCtx("member@example.com"), req: req, expectedCode: codes.PermissionDenied,
		},
		{
			desc: "enforce: not a proto message", mode: Enforce, method: getFeature, ctx: userCtx("viewer@example.com"),
			req: "request", expectedCode: codes.Internal,
		},
		{
			desc: "enforce: health check is skipped", mode: Enforce, method: healthCheck, ctx: context.Background(), req: req,
			expectedHandled: true,
		},
		{
			desc: "enforce: a service not listed stays in shadow", mode: Enforce, method: listNotifications,
			ctx: context.Background(), req: req, expectedHandled: true,
			expectedLogs: []loggedEntry{{level: zapcore.WarnLevel, message: "Would deny"}},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			interceptor, logs := newTestInterceptor(t, p.mode)
			handled := false
			var editor *eventproto.Editor
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				handled = true
				editor = EditorFrom(ctx)
				return "response", p.handlerErr
			}
			resp, err := interceptor.Unary()(p.ctx, p.req, &grpc.UnaryServerInfo{FullMethod: p.method}, handler)
			assert.Equal(t, p.expectedCode, status.Code(err))
			assert.Equal(t, p.expectedHandled, handled)
			if p.expectedHandled {
				assert.Equal(t, "response", resp)
			}
			assert.Equal(t, p.expectedEditor, editor)
			assert.Equal(t, p.expectedLogs, entries(logs))
		})
	}
}

// fakeStream hands the handler the given messages and then io.EOF.
type fakeStream struct {
	grpc.ServerStream
	ctx  context.Context
	msgs []proto.Message
}

func (s *fakeStream) Context() context.Context { return s.ctx }

func (s *fakeStream) RecvMsg(m interface{}) error {
	if len(s.msgs) == 0 {
		return io.EOF
	}
	proto.Merge(m.(proto.Message), s.msgs[0])
	s.msgs = s.msgs[1:]
	return nil
}

func TestStreamInterceptor(t *testing.T) {
	t.Parallel()
	req := &featureproto.GetFeatureRequest{EnvironmentId: "env-1"}
	patterns := []struct {
		desc           string
		mode           Mode
		method         string
		ctx            context.Context
		msgs           []proto.Message
		handlerErr     error
		expectedRecv   int
		expectedEditor *eventproto.Editor
		expectedCode   codes.Code
		expectedLogs   []loggedEntry
	}{
		{
			desc: "shadow: both allow", mode: Shadow, method: getFeature, ctx: userCtx("viewer@example.com"),
			msgs: []proto.Message{req, req}, expectedRecv: 2,
		},
		{
			desc: "shadow: would deny", mode: Shadow, method: getFeature, ctx: userCtx("member@example.com"),
			msgs: []proto.Message{req}, expectedRecv: 1,
			expectedLogs: []loggedEntry{{level: zapcore.WarnLevel, message: "Would deny"}},
		},
		{
			desc: "shadow: would allow", mode: Shadow, method: getFeature, ctx: userCtx("viewer@example.com"),
			msgs: []proto.Message{req}, handlerErr: errDenied, expectedRecv: 1, expectedCode: codes.PermissionDenied,
			expectedLogs: []loggedEntry{{level: zapcore.WarnLevel, message: "Would allow"}},
		},
		{
			desc: "shadow: no message received", mode: Shadow, method: getFeature, ctx: userCtx("member@example.com"),
			expectedRecv: 0,
		},
		{
			desc: "shadow: health check is skipped", mode: Shadow, method: healthCheck, ctx: context.Background(),
			msgs: []proto.Message{req}, expectedRecv: 1,
		},
		{
			desc: "enforce: allow", mode: Enforce, method: getFeature, ctx: userCtx("viewer@example.com"),
			msgs: []proto.Message{req}, expectedRecv: 1, expectedEditor: memberEditor("viewer@example.com"),
		},
		{
			desc: "enforce: deny", mode: Enforce, method: getFeature, ctx: userCtx("member@example.com"),
			msgs: []proto.Message{req}, handlerErr: errNotFound, expectedRecv: 0, expectedCode: codes.PermissionDenied,
		},
		{
			desc: "enforce: health check is skipped", mode: Enforce, method: healthCheck, ctx: context.Background(),
			msgs: []proto.Message{req}, expectedRecv: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			interceptor, logs := newTestInterceptor(t, p.mode)
			received := 0
			var editor *eventproto.Editor
			handler := func(_ interface{}, ss grpc.ServerStream) error {
				for {
					if err := ss.RecvMsg(&featureproto.GetFeatureRequest{}); err != nil {
						if errors.Is(err, io.EOF) {
							break
						}
						return err
					}
					received++
					editor = EditorFrom(ss.Context())
				}
				return p.handlerErr
			}
			stream := &fakeStream{ctx: p.ctx, msgs: p.msgs}
			err := interceptor.Stream()(nil, stream, &grpc.StreamServerInfo{FullMethod: p.method}, handler)
			assert.Equal(t, p.expectedCode, status.Code(err))
			assert.Equal(t, p.expectedRecv, received)
			assert.Equal(t, p.expectedEditor, editor)
			assert.Equal(t, p.expectedLogs, entries(logs))
		})
	}
}

func TestCallerFromContext(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc       string
		ctx        context.Context
		expected   Caller
		expectedOK bool
	}{
		{desc: "no token", ctx: context.Background()},
		{desc: "user", ctx: userCtx("viewer@example.com"), expected: user("viewer@example.com"), expectedOK: true},
		{desc: "service", ctx: serviceCtx(""), expected: Caller{Email: "service@example.com", Service: LegacyService}, expectedOK: true},
		{desc: "relayed", ctx: serviceCtx("viewer@example.com"), expected: relayed("viewer@example.com"), expectedOK: true},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, ok := CallerFromContext(p.ctx)
			assert.Equal(t, p.expectedOK, ok)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestSkipped(t *testing.T) {
	t.Parallel()
	assert.True(t, skipped(healthCheck))
	assert.False(t, skipped(getFeature))
}
