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
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

// HandlerRegistrar is a function that registers a gRPC-Gateway handler
type HandlerRegistrar func(context.Context, *runtime.ServeMux, []grpc.DialOption) error

// customJSONPb is a custom JSONPb that handles both boolean and string representations of boolean fields
type customJSONPb struct {
	runtime.JSONPb
}

// Unmarshal implements the custom unmarshaling logic
// TODO: This is a temporary solution until the Android SDK is updated to use the correct boolean type
func (c *customJSONPb) Unmarshal(data []byte, v interface{}) error {
	// First try to unmarshal with the default protojson unmarshaler
	err := c.JSONPb.Unmarshal(data, v)
	if err == nil {
		return nil
	}

	// Save the original error for later use
	originalErr := err

	// Try to handle boolean strings at the top level
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		// Check if it's a boolean string
		switch strVal {
		case "true", "True", "TRUE", "1":
			return c.JSONPb.Unmarshal([]byte("true"), v)
		case "false", "False", "FALSE", "0":
			return c.JSONPb.Unmarshal([]byte("false"), v)
		default:
			// Not a boolean string, return the original error
			return originalErr
		}
	}

	// Try to handle nested structures with boolean strings
	var rawData interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		// Can't parse as JSON at all, return the original protobuf error
		return originalErr
	}

	// Convert boolean strings in the data structure
	if modified := convertBooleanStrings(rawData); modified {
		// Re-marshal the modified data
		modifiedData, err := json.Marshal(rawData)
		if err != nil {
			return originalErr
		}

		// Try to unmarshal again with the modified data
		if err := c.JSONPb.Unmarshal(modifiedData, v); err == nil {
			return nil
		}
	}

	// Return the original error if nothing worked
	return originalErr
}

// convertBooleanStrings recursively converts string representations of booleans to actual booleans
// Returns true if any modifications were made
func convertBooleanStrings(data interface{}) bool {
	modified := false

	switch v := data.(type) {
	case map[string]interface{}:
		// Handle JSON objects
		for key, value := range v {
			if strVal, ok := value.(string); ok {
				switch strVal {
				case "true", "True", "TRUE", "1":
					v[key] = true
					modified = true
				case "false", "False", "FALSE", "0":
					v[key] = false
					modified = true
				}
			} else {
				// Recursively process nested structures
				if convertBooleanStrings(value) {
					modified = true
				}
			}
		}
	case []interface{}:
		// Handle JSON arrays
		for i, value := range v {
			if strVal, ok := value.(string); ok {
				switch strVal {
				case "true", "True", "TRUE", "1":
					v[i] = true
					modified = true
				case "false", "False", "FALSE", "0":
					v[i] = false
					modified = true
				}
			} else {
				// Recursively process nested structures
				if convertBooleanStrings(value) {
					modified = true
				}
			}
		}
	}

	return modified
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
		logger:   options.logger.Named("gateway"),
	}, nil
}

func (g *Gateway) Start(ctx context.Context,
	registerFuncs ...HandlerRegistrar,
) error {
	// Create a new ServeMux for the REST gateway
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &customJSONPb{
			JSONPb: runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
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
		Addr:    g.restAddr,
		Handler: mux,
	}

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		var err error
		if g.opts.keyPath != "" && g.opts.certPath != "" {
			err = g.httpServer.ListenAndServeTLS(g.opts.certPath, g.opts.keyPath)
		} else {
			g.logger.Info("starting gateway with HTTP (no TLS)")
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
