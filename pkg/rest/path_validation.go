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

package rest

import (
	"net/http"
	"strings"
	"unicode/utf8"
)

// rejectInvalidRequestPathMiddleware rejects malformed or hostile probe paths
// before they reach application handlers (e.g. dashboard static file server).
func rejectInvalidRequestPathMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status, reject := invalidRequestPathStatus(r); reject {
			http.Error(w, http.StatusText(status), status)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func invalidRequestPathStatus(r *http.Request) (int, bool) {
	if r.URL == nil {
		return http.StatusBadRequest, true
	}
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	if hasMalformedPercentEncoding(r.RequestURI) {
		return http.StatusBadRequest, true
	}
	if r.URL.RawPath != "" && hasMalformedPercentEncoding(r.URL.RawPath) {
		return http.StatusBadRequest, true
	}
	if !utf8.ValidString(path) {
		return http.StatusBadRequest, true
	}
	if hasDisallowedControlChars(path) {
		return http.StatusBadRequest, true
	}
	if isCommonSecretProbePath(path) {
		return http.StatusNotFound, true
	}
	return 0, false
}

func hasMalformedPercentEncoding(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '%' {
			continue
		}
		if i+2 >= len(s) {
			return true
		}
		if !isHexDigit(s[i+1]) || !isHexDigit(s[i+2]) {
			return true
		}
		i += 2
	}
	return false
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}

func hasDisallowedControlChars(path string) bool {
	for _, r := range path {
		if r < 0x20 || (r >= 0x7f && r <= 0x9f) {
			return true
		}
	}
	return false
}

func isCommonSecretProbePath(path string) bool {
	lower := strings.ToLower(path)
	for _, prefix := range secretProbePathPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

var secretProbePathPrefixes = []string{
	"/.env",
	"/.git",
	"/.aws/",
	"/.docker/",
	"/credentials",
	"/terraform.tfstate",
	"/azure.yml",
	"/id_rsa",
	"/.ssh/",
	"/wp-admin",
	"/phpmyadmin",
}
