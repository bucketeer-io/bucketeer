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

package rpc

import (
	"context"

	"google.golang.org/grpc"
)

// TODO: change in case this (lambda etc) becomes a performance bottleneck.
func chainUnaryServerInterceptors(is ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// each interceptor should get the context and request from the previous one
	// real grpc handler should be called in the end and passed back up
	// interceptor()
	// -- add some before stuff
	// -- res, err := handler()
	// -- interceptor() <------ this means handler should be interceptor and the last handler should be the real handler
	// ---- add more to context
	// ---- res, err := handler()
	// ---- do some after stuff
	// ---- return res, err
	// -- do some after stuff
	// return res, err
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		chain := func(interceptor grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return interceptor(ctx, req, info, next)
			}
		}
		next := handler
		for i := len(is) - 1; i >= 0; i-- {
			next = chain(is[i], next)
		}
		return next(ctx, req)
	}
}
