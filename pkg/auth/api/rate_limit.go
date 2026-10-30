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

package api

import (
	"context"
	"net"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
)

// xForwardedForKey is the metadata key under which grpc-gateway forwards the
// X-Forwarded-For header, with the remote address of the HTTP client appended.
const xForwardedForKey = "x-forwarded-for"

// unknownClientAddress is used as the rate limit key when no client address
// can be determined from the request context.
const unknownClientAddress = "unknown"

// throttle checks whether the caller identified by key may run the given
// operation. A nil limiter disables throttling. When the limit is exceeded,
// the denial is logged and statusTooManyRequests is returned.
func (s *authService) throttle(
	ctx context.Context,
	limiter *ratelimit.Limiter,
	operation string,
	key string,
) error {
	if limiter == nil {
		return nil
	}
	if limiter.Allow(operation + ":" + key) {
		return nil
	}
	s.logger.Warn(
		"Too many requests",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.String("operation", operation),
			zap.String("key", key),
		)...,
	)
	return statusTooManyRequests.Err()
}

// clientAddressFromContext returns the address of the client that sent the
// request. It prefers the first address of the x-forwarded-for metadata,
// which grpc-gateway always populates for requests received over HTTP, and
// falls back to the gRPC peer address.
func clientAddressFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for _, value := range md.Get(xForwardedForKey) {
			first, _, _ := strings.Cut(value, ",")
			if addr := strings.TrimSpace(first); addr != "" {
				return addr
			}
		}
	}
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		addr := p.Addr.String()
		if host, _, err := net.SplitHostPort(addr); err == nil {
			return host
		}
		if addr != "" {
			return addr
		}
	}
	return unknownClientAddress
}
