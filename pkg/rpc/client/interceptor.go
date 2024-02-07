// Copyright 2024 The Bucketeer Authors.
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

package client

import (
	"context"

	"google.golang.org/grpc"
)

func ChainUnaryClientInterceptors(is ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		chain := func(interceptor grpc.UnaryClientInterceptor, next grpc.UnaryInvoker) grpc.UnaryInvoker {
			return func(
				ctx context.Context,
				method string,
				req, reply interface{},
				cc *grpc.ClientConn,
				opts ...grpc.CallOption,
			) error {
				return interceptor(ctx, method, req, reply, cc, next, opts...)
			}
		}
		next := invoker
		for i := len(is) - 1; i >= 0; i-- {
			next = chain(is[i], next)
		}
		return next(ctx, method, req, reply, cc, opts...)
	}
}
