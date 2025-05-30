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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// BooleanConversionMiddleware handles conversion of string boolean values to actual booleans
// This is specifically needed for SDKs that send "true"/"false" strings instead of boolean values
func BooleanConversionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only process POST requests with JSON content
		if r.Method != http.MethodPost || !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			next.ServeHTTP(w, r)
			return
		}

		// Only process specific endpoints that have boolean conversion issues
		if !shouldProcessRequest(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// Convert string booleans to actual booleans
		convertedBody, err := convertStringBooleans(body)
		if err != nil {
			// If conversion fails, pass through the original body
			convertedBody = body
		}

		// Create a new request with the converted body
		r.Body = io.NopCloser(bytes.NewReader(convertedBody))
		r.ContentLength = int64(len(convertedBody))

		next.ServeHTTP(w, r)
	})
}

// shouldProcessRequest determines if the request should be processed for boolean conversion
func shouldProcessRequest(path string) bool {
	// Add paths that need boolean conversion
	conversionPaths := []string{
		"/get_evaluations",
		// Add other paths as needed
	}

	for _, conversionPath := range conversionPaths {
		if path == conversionPath {
			return true
		}
	}
	return false
}

// convertStringBooleans converts string boolean values to actual booleans in JSON
func convertStringBooleans(body []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body, err
	}

	// Convert string booleans recursively
	converted := convertStringBooleansRecursive(data)

	// Marshal back to JSON
	result, err := json.Marshal(converted)
	if err != nil {
		return body, err
	}

	return result, nil
}

// convertStringBooleansRecursive recursively converts string boolean values
func convertStringBooleansRecursive(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			// Convert specific fields that should be booleans
			if shouldConvertField(key) {
				if strVal, ok := value.(string); ok {
					if boolVal, converted := stringToBool(strVal); converted {
						result[key] = boolVal
					} else {
						result[key] = value
					}
				} else {
					result[key] = convertStringBooleansRecursive(value)
				}
			} else {
				result[key] = convertStringBooleansRecursive(value)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = convertStringBooleansRecursive(item)
		}
		return result
	default:
		return data
	}
}

// shouldConvertField determines if a field should be converted from string to boolean
func shouldConvertField(fieldName string) bool {
	// List of fields that should be converted from string to boolean
	booleanFields := []string{
		"userAttributesUpdated",
		"user_attributes_updated",
		// Add other boolean fields as needed
	}

	for _, field := range booleanFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// stringToBool converts string boolean values to actual booleans
func stringToBool(s string) (bool, bool) {
	switch strings.ToLower(s) {
	case "true", "1":
		return true, true
	case "false", "0":
		return false, true
	default:
		return false, false
	}
}
