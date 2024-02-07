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
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb" // nolint:staticcheck
	"github.com/golang/protobuf/proto"  // nolint:staticcheck
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

var marshaler = jsonpb.Marshaler{EmitDefaults: true}

func LogUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	logger = logger.Named("grpc_server")
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		level := zap.DebugLevel
		code := status.Code(err)
		if err != nil {
			switch code {
			case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated:
				level = zap.WarnLevel
			default:
				level = zap.ErrorLevel
			}
		}
		if level == zap.DebugLevel && strings.HasPrefix(info.FullMethod, healthServiceName) {
			return resp, err
		}
		serviceName, methodName := splitFullMethodName(info.FullMethod)
		logger.Check(level, "").Write(
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("grpcService", serviceName),
				zap.String("grpcMethod", methodName),
				zap.String("grpcCode", code.String()),
				zap.Duration("duration", time.Since(startTime)),
				zap.Reflect("request", makeMarshallable(req)),
				zap.Reflect("response", makeMarshallable(resp)),
			)...,
		)
		return resp, err
	}
}

type marshallable struct {
	proto.Message
}

func makeMarshallable(msg interface{}) *marshallable {
	if m, ok := msg.(proto.Message); ok {
		return &marshallable{m}
	}
	return nil
}

/*
zap json encoder calls json.Marshal when doing zap.Reflected()

	func (enc *jsonEncoder) AppendReflected(val interface{}) error {
		marshaled, err := json.Marshal(val)
		if err != nil {
			return err
		}
		enc.addElementSeparator()
		_, err = enc.buf.Write(marshaled)
		return err
	}
*/
func (m *marshallable) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	if err := marshaler.Marshal(b, m.Message); err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}
	return b.Bytes(), nil
}
