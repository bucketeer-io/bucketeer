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

package gateway

import (
	"net/textproto"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// AllowedIncomingHeaderMatcher returns a runtime.HeaderMatcherFunc that
// forwards only the given HTTP request headers to the gRPC services, as
// metadata keyed by the lower-cased header name. Every other header is
// dropped.
//
// Unlike runtime.DefaultHeaderMatcher, it never maps "Grpc-Metadata-*"
// headers, so HTTP clients cannot set gRPC metadata that is reserved for
// internal callers.
func AllowedIncomingHeaderMatcher(headers ...string) runtime.HeaderMatcherFunc {
	allowed := make(map[string]string, len(headers))
	for _, header := range headers {
		allowed[textproto.CanonicalMIMEHeaderKey(header)] = strings.ToLower(header)
	}
	return func(key string) (string, bool) {
		metadataKey, ok := allowed[textproto.CanonicalMIMEHeaderKey(key)]
		return metadataKey, ok
	}
}
