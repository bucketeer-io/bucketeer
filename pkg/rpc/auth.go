// Copyright 2023 The Bucketeer Authors.
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

package rpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/token"
)

type tokenKey struct{}

var Key = tokenKey{}

const (
	healthServiceName      = "/grpc.health.v1.Health/"
	flagTriggerWebhookName = "/bucketeer.feature.FeatureService/FlagTriggerWebhook"
)

func AuthUnaryServerInterceptor(verifier token.Verifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.HasPrefix(info.FullMethod, healthServiceName) ||
			strings.HasPrefix(info.FullMethod, flagTriggerWebhookName) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "token is required")
		}
		rawTokens, ok := md["authorization"]
		if !ok || len(rawTokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token is required")
		}
		subs := strings.Split(rawTokens[0], " ")
		if len(subs) != 2 {
			return nil, status.Error(codes.Unauthenticated, "token is malformed")
		}
		token, err := verifier.Verify(subs[1])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %s", err.Error())
		}
		ctx = context.WithValue(ctx, Key, token)
		return handler(ctx, req)
	}
}

func GetIDToken(ctx context.Context) (*token.IDToken, bool) {
	t, ok := ctx.Value(Key).(*token.IDToken)
	return t, ok
}
