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
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBooleanConversionMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		method        string
		contentType   string
		requestBody   string
		expectedBody  string
		shouldConvert bool
	}{
		{
			name:        "Convert string boolean to actual boolean",
			path:        "/get_evaluations",
			method:      "POST",
			contentType: "application/json",
			requestBody: `{
				"tag": "production",
				"userEvaluationCondition": {
					"evaluatedAt": "1234567890",
					"userAttributesUpdated": "false"
				}
			}`,
			expectedBody: `{
				"tag": "production",
				"userEvaluationCondition": {
					"evaluatedAt": "1234567890",
					"userAttributesUpdated": false
				}
			}`,
			shouldConvert: true,
		},
		{
			name:        "Convert snake_case field",
			path:        "/get_evaluations",
			method:      "POST",
			contentType: "application/json",
			requestBody: `{
				"tag": "production",
				"userEvaluationCondition": {
					"evaluatedAt": "1234567890",
					"user_attributes_updated": "true"
				}
			}`,
			expectedBody: `{
				"tag": "production",
				"userEvaluationCondition": {
					"evaluatedAt": "1234567890",
					"user_attributes_updated": true
				}
			}`,
			shouldConvert: true,
		},
		{
			name:        "Don't convert non-boolean fields",
			path:        "/get_evaluations",
			method:      "POST",
			contentType: "application/json",
			requestBody: `{
				"tag": "production",
				"someOtherField": "false"
			}`,
			expectedBody: `{
				"tag": "production",
				"someOtherField": "false"
			}`,
			shouldConvert: true,
		},
		{
			name:        "Skip non-target endpoints",
			path:        "/other_endpoint",
			method:      "POST",
			contentType: "application/json",
			requestBody: `{
				"userAttributesUpdated": "false"
			}`,
			expectedBody: `{
				"userAttributesUpdated": "false"
			}`,
			shouldConvert: false,
		},
		{
			name:        "Skip GET requests",
			path:        "/get_evaluations",
			method:      "GET",
			contentType: "application/json",
			requestBody: `{
				"userAttributesUpdated": "false"
			}`,
			expectedBody: `{
				"userAttributesUpdated": "false"
			}`,
			shouldConvert: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that captures the request body
			var capturedBody []byte
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				capturedBody = body
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with our middleware
			middleware := BooleanConversionMiddleware(testHandler)

			// Create test request
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", tt.contentType)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute the request
			middleware.ServeHTTP(rr, req)

			// Parse the captured body
			var capturedJSON interface{}
			if err := json.Unmarshal(capturedBody, &capturedJSON); err != nil {
				t.Fatalf("Failed to parse captured body: %v", err)
			}

			// Parse the expected body
			var expectedJSON interface{}
			if err := json.Unmarshal([]byte(tt.expectedBody), &expectedJSON); err != nil {
				t.Fatalf("Failed to parse expected body: %v", err)
			}

			// Compare the JSON structures
			capturedBytes, _ := json.Marshal(capturedJSON)
			expectedBytes, _ := json.Marshal(expectedJSON)

			if !bytes.Equal(capturedBytes, expectedBytes) {
				t.Errorf("Body mismatch.\nExpected: %s\nGot: %s", string(expectedBytes), string(capturedBytes))
			}
		})
	}
}

func TestStringToBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		valid    bool
	}{
		{"true", true, true},
		{"True", true, true},
		{"TRUE", true, true},
		{"1", true, true},
		{"false", false, true},
		{"False", false, true},
		{"FALSE", false, true},
		{"0", false, true},
		{"invalid", false, false},
		{"", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, valid := stringToBool(tt.input)
			if result != tt.expected || valid != tt.valid {
				t.Errorf("stringToBool(%q) = (%v, %v), want (%v, %v)",
					tt.input, result, valid, tt.expected, tt.valid)
			}
		})
	}
}

func TestShouldConvertField(t *testing.T) {
	tests := []struct {
		fieldName string
		expected  bool
	}{
		{"userAttributesUpdated", true},
		{"user_attributes_updated", true},
		{"someOtherField", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			result := shouldConvertField(tt.fieldName)
			if result != tt.expected {
				t.Errorf("shouldConvertField(%q) = %v, want %v",
					tt.fieldName, result, tt.expected)
			}
		})
	}
}
