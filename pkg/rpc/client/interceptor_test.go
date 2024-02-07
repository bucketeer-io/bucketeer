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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestChainUnaryClientInterceptors(t *testing.T) {
	type parentKey string
	parent := parentKey("parent")
	ctx := context.WithValue(context.Background(), parent, "")
	serviceMethod := "test-method"
	var firstRun, secondRun, invokerRun bool
	first := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		require.Equal(t, serviceMethod, method)
		require.Equal(t, "", ctx.Value(parent).(string))
		ctx = context.WithValue(ctx, parent, "first")
		firstRun = true
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	second := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		require.Equal(t, serviceMethod, method)
		require.Equal(t, "first", ctx.Value(parent).(string))
		ctx = context.WithValue(ctx, parent, "second")
		secondRun = true
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		require.Equal(t, serviceMethod, method)
		require.Equal(t, "second", ctx.Value(parent).(string))
		invokerRun = true
		return nil
	}
	interceptors := ChainUnaryClientInterceptors(first, second)
	interceptors(ctx, serviceMethod, "req", "reply", nil, invoker, nil)
	assert.True(t, firstRun)
	assert.True(t, secondRun)
	assert.True(t, invokerRun)
}
