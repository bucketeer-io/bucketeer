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

package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

func LogServerMiddleware(logger *zap.Logger) middleware {
	return middleware(
		func(next http.Handler) http.Handler {
			logger = logger.Named("http_server")
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					startTime := time.Now()
					rr := &responseRecorder{
						ResponseWriter: w,
						statusCode:     200,
						body:           new(bytes.Buffer),
					}
					next.ServeHTTP(rr, r)
					if rr.statusCode == http.StatusOK {
						return
					}
					var level zapcore.Level
					switch rr.statusCode {
					case http.StatusBadRequest, http.StatusNotFound, http.StatusUnauthorized:
						level = zap.WarnLevel
					default:
						level = zap.ErrorLevel
					}
					apiVersion, serviceName, apiName := splitURLPath(r.URL.Path)
					reqBody, err := decodeBody(r.Body)
					if err != nil {
						logger.Error("Failed to parse request body", zap.Error(err))
					}
					logger.Check(level, "").Write(
						log.FieldsFromImcomingContext(r.Context()).AddFields(
							zap.String("requestURI", r.RequestURI),
							zap.String("apiVersion", apiVersion),
							zap.String("serviceName", serviceName),
							zap.String("apiName", apiName),
							zap.String("httpMethod", r.Method),
							zap.Int("statusCode", rr.statusCode),
							zap.Duration("duration", time.Since(startTime)),
							zap.Any("request", reqBody),
							zap.String("response", rr.body.String()),
						)...,
					)
				},
			)
		},
	)
}

func decodeBody(body io.Reader) (interface{}, error) {
	var decoded interface{}
	if err := json.NewDecoder(body).Decode(&decoded); err != nil {
		if err == io.EOF {
			return decoded, nil
		}
		return nil, err
	}
	return decoded, nil
}
