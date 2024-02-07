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
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb" // nolint:staticcheck
	"github.com/golang/protobuf/proto"  // nolint:staticcheck
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

var marshaler = &jsonpb.Marshaler{}

func LogUnaryClientInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	logger = logger.Named("grpc_client")
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		serviceName, methodName := splitFullMethodName(method)
		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		logger.Check(zap.DebugLevel, "").Write(
			log.FieldsFromOutgoingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("gprcService", serviceName),
				zap.String("grpcMethod", methodName),
				zap.String("grpcCode", status.Code(err).String()),
				zap.Duration("duration", time.Since(startTime)),
				zap.Reflect("request", makeUnmarshallable(req)),
			)...,
		)
		return err
	}
}

type unmarshallable struct {
	proto.Message
}

func makeUnmarshallable(msg interface{}) *unmarshallable {
	if m, ok := msg.(proto.Message); ok {
		return &unmarshallable{m}
	}
	return nil
}

func (m *unmarshallable) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	if err := marshaler.Marshal(b, m); err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}
	return b.Bytes(), nil
}
