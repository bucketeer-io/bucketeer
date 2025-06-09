// Copyright 2025 The Bucketeer Authors.
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

package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

// HandlerRegistrar is a function that registers a gRPC-Gateway handler
type HandlerRegistrar func(context.Context, *runtime.ServeMux, []grpc.DialOption) error

var (
	// Pre-compiled maps for O(1) lookups instead of O(n) iterations
	booleanFieldsMap = map[string]bool{
		"userAttributesUpdated":   true,
		"user_attributes_updated": true,
	}

	// Pre-compiled boolean string values for faster lookups
	stringBoolMap = map[string]bool{
		"true":  true,
		"false": false,
		"1":     true,
		"0":     false,
	}
)

// BooleanConversionMarshaler is a custom marshaler that handles boolean string conversion
// during the unmarshaling phase. This is a temporary solution to handle boolean fields
// sent as strings from the Android SDK.
// TODO: Remove this once the Android SDK is fixed to send proper boolean values.
// Reference: https://github.com/bucketeer-io/android-client-sdk/pull/230
type BooleanConversionMarshaler struct {
	runtime.JSONPb
	logger *zap.Logger
}

// Unmarshal implements the Marshaler interface with boolean string conversion
func (m *BooleanConversionMarshaler) Unmarshal(data []byte, v interface{}) error {
	// First, try to unmarshal as a proto message
	if msg, ok := v.(proto.Message); ok {
		// Check if we need to convert boolean strings based on the message type
		// For now, we'll do a pre-processing step on the JSON data
		convertedData := m.preprocessJSON(data)

		// Use the standard protojson unmarshaler
		err := protojson.UnmarshalOptions{
			DiscardUnknown: m.UnmarshalOptions.DiscardUnknown,
			AllowPartial:   m.UnmarshalOptions.AllowPartial,
		}.Unmarshal(convertedData, msg)

		if err != nil {
			// TODO: Send metrics for 404 errors instead of logging
			// Removed logging to reduce log volume from automated scanners
			return err
		}

		return nil
	}

	// Fall back to standard JSON unmarshaling for non-proto messages
	err := json.Unmarshal(data, v)
	if err != nil {
		m.logger.Error("Failed to unmarshal non-proto message",
			zap.Error(err),
			zap.String("data", string(data)),
		)
	}
	return err
}

// preprocessJSON converts string boolean values to actual booleans in the JSON data
func (m *BooleanConversionMarshaler) preprocessJSON(data []byte) []byte {
	// Fast path: empty data
	if len(data) == 0 {
		return data
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		// Log the error before returning original data
		m.logger.Error("Failed to unmarshal JSON data in preprocessJSON",
			zap.Error(err),
			zap.String("data", string(data)),
		)
		// If we can't unmarshal, return the original data
		return data
	}

	// Convert string booleans recursively
	converted, modified := m.convertStringBooleansRecursive(jsonData)
	if !modified {
		return data
	}

	// Re-marshal the converted data
	convertedData, err := json.Marshal(converted)
	if err != nil {
		m.logger.Error("Failed to re-marshal converted data",
			zap.Error(err),
			zap.Any("converted", converted),
		)
		return data
	}

	return convertedData
}

// convertStringBooleansRecursive recursively converts string boolean values
func (m *BooleanConversionMarshaler) convertStringBooleansRecursive(data interface{}) (interface{}, bool) {
	modified := false

	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))

		for key, value := range v {
			// Check if this field should be converted
			if booleanFieldsMap[key] {
				if strVal, ok := value.(string); ok {
					if boolVal, converted := stringToBool(strVal); converted {
						result[key] = boolVal
						modified = true
					} else {
						result[key] = value
						// Log failed conversion attempt
						m.logger.Warn("Failed to convert boolean field",
							zap.String("field", key),
							zap.String("value", strVal),
						)
					}
				} else {
					convertedValue, childModified := m.convertStringBooleansRecursive(value)
					result[key] = convertedValue
					modified = modified || childModified
				}
			} else {
				convertedValue, childModified := m.convertStringBooleansRecursive(value)
				result[key] = convertedValue
				modified = modified || childModified
			}
		}
		return result, modified

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			convertedItem, childModified := m.convertStringBooleansRecursive(item)
			result[i] = convertedItem
			modified = modified || childModified
		}
		return result, modified

	default:
		return data, false
	}
}

// stringToBool converts string boolean values to actual booleans
func stringToBool(s string) (bool, bool) {
	if val, ok := stringBoolMap[strings.ToLower(strings.TrimSpace(s))]; ok {
		return val, true
	}
	return false, false
}

// NewDecoder returns a decoder that handles boolean string conversion
func (m *BooleanConversionMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &booleanConversionDecoder{
		reader:    r,
		marshaler: m,
	}
}

// booleanConversionDecoder is a custom decoder that handles boolean string conversion
type booleanConversionDecoder struct {
	reader    io.Reader
	marshaler *BooleanConversionMarshaler
}

// Decode implements the Decoder interface
func (d *booleanConversionDecoder) Decode(v interface{}) error {
	data, err := io.ReadAll(d.reader)
	if err != nil {
		// TODO: Send metrics for read errors instead of logging to avoid log volume
		return err
	}
	return d.marshaler.Unmarshal(data, v)
}

// noOpWriter discards all logs from the HTTP server
type noOpWriter struct{}

func (noOpWriter) Write(p []byte) (n int, err error) {
	// Discard all logs
	return len(p), nil
}

type Gateway struct {
	httpServer *http.Server
	restAddr   string
	opts       *options
	logger     *zap.Logger
}

func NewGateway(restAddr string, opts ...Option) (*Gateway, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.logger == nil {
		logger, err := log.NewLogger()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %v", err)
		}
		options.logger = logger
	}

	return &Gateway{
		restAddr: restAddr,
		opts:     &options,
		logger:   options.logger.Named("grpc-gateway"),
	}, nil
}

// customRoutingErrorHandler handles routing errors and strips trailing slashes
func (g *Gateway) customRoutingErrorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	httpStatus int,
) {
	// If we get a 404 and the path ends with a slash, try without the slash
	if httpStatus == http.StatusNotFound && strings.HasSuffix(r.URL.Path, "/") {
		// Create a new request with the trailing slash removed
		newPath := strings.TrimSuffix(r.URL.Path, "/")
		// Clone the request with the new path
		newReq := r.Clone(ctx)
		newReq.URL.Path = newPath
		newReq.RequestURI = newPath
		if r.URL.RawQuery != "" {
			newReq.RequestURI += "?" + r.URL.RawQuery
		}

		g.logger.Debug("retrying request without trailing slash",
			zap.String("original_full_url", r.URL.String()),
			zap.String("new_full_url", newReq.URL.String()),
		)

		// Try to serve the request again
		mux.ServeHTTP(w, newReq)
		return
	}

	// For other cases, use the default routing error handler
	runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
}

func (g *Gateway) Start(ctx context.Context,
	registerFuncs ...HandlerRegistrar,
) error {
	// Create our custom marshaler that handles boolean conversion
	customMarshaler := &BooleanConversionMarshaler{
		JSONPb: runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
		logger: g.logger,
	}

	// Custom error handler that logs errors
	errorHandler := func(
		ctx context.Context,
		mux *runtime.ServeMux,
		marshaler runtime.Marshaler,
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {
		// TODO: Send metrics for 404 errors instead of logging
		// Removed logging to reduce log volume from automated scanners

		// Call the default error handler to send the response
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
	}

	// Create a new ServeMux with our custom marshaler
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(errorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, customMarshaler),
		runtime.WithRoutingErrorHandler(g.customRoutingErrorHandler),
	)

	// Create gRPC dial options with proper credentials and settings
	var dialOpts []grpc.DialOption
	if g.opts.certPath != "" {
		creds, err := credentials.NewClientTLSFromFile(g.opts.certPath, "")
		if err != nil {
			return fmt.Errorf("failed to create TLS credentials: %v", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		g.logger.Warn("starting gateway without TLS credentials")
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add keepalive parameters
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                g.opts.keepaliveTime,
		Timeout:             g.opts.keepaliveTimeout,
		PermitWithoutStream: g.opts.permitWithoutStream,
	}))

	// Add window size parameters
	dialOpts = append(dialOpts,
		grpc.WithInitialWindowSize(g.opts.initialWindowSize),
		grpc.WithInitialConnWindowSize(g.opts.initialConnWindowSize),
	)

	// Register all the provided handler registrars
	for _, registerFunc := range registerFuncs {
		if err := registerFunc(ctx, mux, dialOpts); err != nil {
			return fmt.Errorf("failed to register gateway handler: %v", err)
		}
	}

	// Create and start the HTTP server
	g.httpServer = &http.Server{
		Addr:     g.restAddr,
		Handler:  mux,
		ErrorLog: stdlog.New(noOpWriter{}, "", 0),
	}

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		var err error
		if g.opts.keyPath != "" && g.opts.certPath != "" {
			err = g.httpServer.ListenAndServeTLS(g.opts.certPath, g.opts.keyPath)
		} else {
			err = g.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Check if there was an immediate error
	select {
	case err := <-errChan:
		return fmt.Errorf("failed to start gateway: %v", err)
	default:
		// No immediate error, server is starting
		g.logger.Debug("gateway started",
			zap.String("rest_addr", g.restAddr),
			zap.Bool("tls_enabled", g.opts.keyPath != "" && g.opts.certPath != ""),
			zap.Int("handlers_registered", len(registerFuncs)))
		return nil
	}
}

// Stop gracefully shuts down the HTTP server
func (g *Gateway) Stop(timeout time.Duration) {
	if g.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := g.httpServer.Shutdown(shutdownCtx); err != nil {
			g.logger.Error("failed to shutdown HTTP server gracefully", zap.Error(err))
		}
	}
}
