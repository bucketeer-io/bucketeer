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

package oidc

import (
	"testing"
)

func TestGenerateNonce(t *testing.T) {
	nonce, err := GenerateNonce()
	if err != nil {
		t.Fatalf("GenerateNonce() error = %v", err)
	}

	if nonce == "" {
		t.Error("GenerateNonce() returned empty nonce")
	}

	// Verify nonce can be validated
	if err := ValidateNonce(nonce); err != nil {
		t.Errorf("GenerateNonce() produced invalid nonce: %v", err)
	}
}

func TestGenerateNonceUniqueness(t *testing.T) {
	nonce1, err := GenerateNonce()
	if err != nil {
		t.Fatalf("GenerateNonce() error = %v", err)
	}

	nonce2, err := GenerateNonce()
	if err != nil {
		t.Fatalf("GenerateNonce() error = %v", err)
	}

	if nonce1 == nonce2 {
		t.Error("GenerateNonce() generated same nonce twice")
	}
}

func TestValidateNonce(t *testing.T) {
	tests := []struct {
		name      string
		nonce     string
		wantError bool
	}{
		{
			name:      "valid nonce",
			nonce:     "dGVzdC1ub25jZQ",
			wantError: false,
		},
		{
			name:      "empty nonce",
			nonce:     "",
			wantError: true,
		},
		{
			name:      "invalid base64",
			nonce:     "not-valid-base64!!!",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonce(tt.nonce)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNonce() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
