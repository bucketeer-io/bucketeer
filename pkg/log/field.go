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

package log

import (
	"context"
	"errors"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bucketeer-io/bucketeer/pkg/rpc/metadata"
)

type Fields []zap.Field

func FieldsFromImcomingContext(ctx context.Context) Fields {
	sc := trace.FromContext(ctx).SpanContext()
	return Fields{
		zap.String("xRequestID", metadata.GetXRequestIDFromIncomingContext(ctx)),
		zap.String("logging.googleapis.com/trace", sc.TraceID.String()),
		zap.String("logging.googleapis.com/spanId", sc.SpanID.String()),
	}
}

func FieldsFromOutgoingContext(ctx context.Context) Fields {
	sc := trace.FromContext(ctx).SpanContext()
	return Fields{
		zap.String("xRequestID", metadata.GetXRequestIDFromOutgoingContext(ctx)),
		zap.String("logging.googleapis.com/trace", sc.TraceID.String()),
		zap.String("logging.googleapis.com/spanId", sc.SpanID.String()),
	}
}

func (fs Fields) AddFields(fields ...zap.Field) Fields {
	return append(fs, fields...)
}

type serviceContext struct {
	service string
	version string
}

func (sc *serviceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if sc.service == "" {
		return errors.New("service name is mandatory")
	}
	enc.AddString("service", sc.service)
	enc.AddString("version", sc.version)
	return nil
}
