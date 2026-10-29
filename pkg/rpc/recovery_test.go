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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testFullMethod = "/bucketeer.test.TestService/Test"

type fakeServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *fakeServerStream) Context() context.Context {
	return s.ctx
}

// nilPointerPanic dereferences a nil pointer, which is the kind of runtime
// error the interceptors must turn into a response instead of a crash.
func nilPointerPanic() {
	var resp *struct{ message string }
	_ = resp.message
}

func TestRecoveryUnaryServerInterceptor(t *testing.T) {
	t.Parallel()
	handlerErr := status.Error(codes.NotFound, "not found")
	patterns := []struct {
		desc         string
		handler      grpc.UnaryHandler
		expectedResp interface{}
		expectedErr  error
		expectedCode codes.Code
	}{
		{
			desc: "success: response is passed through",
			handler: func(context.Context, interface{}) (interface{}, error) {
				return "response", nil
			},
			expectedResp: "response",
			expectedCode: codes.OK,
		},
		{
			desc: "err: handler error is passed through unchanged",
			handler: func(context.Context, interface{}) (interface{}, error) {
				return nil, handlerErr
			},
			expectedErr:  handlerErr,
			expectedCode: codes.NotFound,
		},
		{
			desc: "err: explicit panic becomes an internal error",
			handler: func(context.Context, interface{}) (interface{}, error) {
				panic("boom")
			},
			expectedCode: codes.Internal,
		},
		{
			desc: "err: nil pointer dereference becomes an internal error",
			handler: func(context.Context, interface{}) (interface{}, error) {
				nilPointerPanic()
				return "unreachable", nil
			},
			expectedCode: codes.Internal,
		},
		{
			desc: "err: panic with an error value becomes an internal error",
			handler: func(context.Context, interface{}) (interface{}, error) {
				panic(errors.New("boom"))
			},
			expectedCode: codes.Internal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			interceptor := RecoveryUnaryServerInterceptor(zap.NewNop())
			info := &grpc.UnaryServerInfo{FullMethod: testFullMethod}
			resp, err := interceptor(context.Background(), "request", info, p.handler)
			assert.Equal(t, p.expectedResp, resp)
			assert.Equal(t, p.expectedCode, status.Code(err))
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			}
			if p.expectedCode == codes.Internal {
				require.Error(t, err)
				// The panic details must not leak to the client.
				assert.Equal(t, "internal error", status.Convert(err).Message())
			}
		})
	}
}

func TestRecoveryStreamServerInterceptor(t *testing.T) {
	t.Parallel()
	handlerErr := status.Error(codes.PermissionDenied, "denied")
	patterns := []struct {
		desc         string
		handler      grpc.StreamHandler
		expectedErr  error
		expectedCode codes.Code
	}{
		{
			desc: "success: nil error is passed through",
			handler: func(interface{}, grpc.ServerStream) error {
				return nil
			},
			expectedCode: codes.OK,
		},
		{
			desc: "err: handler error is passed through unchanged",
			handler: func(interface{}, grpc.ServerStream) error {
				return handlerErr
			},
			expectedErr:  handlerErr,
			expectedCode: codes.PermissionDenied,
		},
		{
			desc: "err: explicit panic becomes an internal error",
			handler: func(interface{}, grpc.ServerStream) error {
				panic("boom")
			},
			expectedCode: codes.Internal,
		},
		{
			desc: "err: nil pointer dereference becomes an internal error",
			handler: func(interface{}, grpc.ServerStream) error {
				nilPointerPanic()
				return nil
			},
			expectedCode: codes.Internal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			interceptor := RecoveryStreamServerInterceptor(zap.NewNop())
			info := &grpc.StreamServerInfo{FullMethod: testFullMethod, IsServerStream: true}
			stream := &fakeServerStream{ctx: context.Background()}
			err := interceptor(nil, stream, info, p.handler)
			assert.Equal(t, p.expectedCode, status.Code(err))
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			}
			if p.expectedCode == codes.Internal {
				require.Error(t, err)
				assert.Equal(t, "internal error", status.Convert(err).Message())
			}
		})
	}
}
