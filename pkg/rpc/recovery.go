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

package rpc

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
)

// RecoveryUnaryServerInterceptor recovers from a panic raised while handling a unary RPC
// and returns a codes.Internal error instead, so a single request cannot terminate the process.
// gRPC runs every handler in its own goroutine, so a panic that is not recovered here crashes
// the whole server. It must be the last interceptor in the chain so the outer interceptors
// (metrics, auth) observe the resulting error like any other handler error.
func RecoveryUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	logger = logger.Named("recovery")
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				resp = nil
				err = recoverFromPanic(ctx, logger, info.FullMethod, r)
			}
		}()
		return handler(ctx, req)
	}
}

// RecoveryStreamServerInterceptor is the streaming counterpart of RecoveryUnaryServerInterceptor.
func RecoveryStreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	logger = logger.Named("recovery")
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFromPanic(ss.Context(), logger, info.FullMethod, r)
			}
		}()
		return handler(srv, ss)
	}
}

func recoverFromPanic(ctx context.Context, logger *zap.Logger, fullMethod string, p interface{}) error {
	serviceName, methodName := splitFullMethodName(fullMethod)
	logger.Error("Recovered from a panic while handling an RPC",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Any("panic", p),
			zap.String("grpcService", serviceName),
			zap.String("grpcMethod", methodName),
			zap.String("stack", string(debug.Stack())),
		)...,
	)
	// Do not expose the panic details to the client.
	return status.Error(codes.Internal, "internal error")
}
