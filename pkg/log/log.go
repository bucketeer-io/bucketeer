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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Levels = []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

type options struct {
	level          string
	serviceContext *serviceContext
}

type Option func(*options)

func WithLevel(level string) Option {
	return func(opts *options) {
		opts.level = level
	}
}

func WithServiceContext(service, version string) Option {
	return func(opts *options) {
		opts.serviceContext = &serviceContext{
			service: service,
			version: version,
		}
	}
}

func NewLogger(opts ...Option) (*zap.Logger, error) {
	dopts := &options{
		level: "info",
	}
	for _, opt := range opts {
		opt(dopts)
	}
	level := new(zapcore.Level)
	if err := level.Set(dopts.level); err != nil {
		return nil, err
	}
	if dopts.serviceContext == nil {
		return newConfig(*level).Build()
	}
	option := zap.Fields(zap.Object("serviceContext", dopts.serviceContext))
	logger, err := newConfig(*level).Build(option)
	if err != nil {
		return nil, err
	}
	return logger.Named(dopts.serviceContext.service), nil
}

func newConfig(level zapcore.Level) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    newEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "eventTime",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARNING")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("CRITICAL")
	case zapcore.PanicLevel:
		enc.AppendString("ALERT")
	case zapcore.FatalLevel:
		enc.AppendString("EMERGENCY")
	}
}
