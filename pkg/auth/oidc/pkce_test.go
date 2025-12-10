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

func TestGeneratePKCEChallenge(t *testing.T) {
	challenge, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge() error = %v", err)
	}

	if challenge.Verifier == "" {
		t.Error("GeneratePKCEChallenge() verifier is empty")
	}

	if challenge.Challenge == "" {
		t.Error("GeneratePKCEChallenge() challenge is empty")
	}

	if challenge.ChallengeMethod != "S256" {
		t.Errorf("GeneratePKCEChallenge() method = %v, want S256", challenge.ChallengeMethod)
	}

	// Verify verifier length is within acceptable range
	if len(challenge.Verifier) < 43 || len(challenge.Verifier) > 128 {
		t.Errorf("GeneratePKCEChallenge() verifier length = %d, want between 43 and 128", len(challenge.Verifier))
	}
}

func TestGeneratePKCEChallengeUniqueness(t *testing.T) {
	challenge1, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge() error = %v", err)
	}

	challenge2, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge() error = %v", err)
	}

	if challenge1.Verifier == challenge2.Verifier {
		t.Error("GeneratePKCEChallenge() generated same verifier twice")
	}

	if challenge1.Challenge == challenge2.Challenge {
		t.Error("GeneratePKCEChallenge() generated same challenge twice")
	}
}

func TestValidateCodeVerifier(t *testing.T) {
	tests := []struct {
		name      string
		verifier  string
		wantError bool
	}{
		{
			name:      "valid verifier",
			verifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
			wantError: false,
		},
		{
			name:      "too short",
			verifier:  "short",
			wantError: true,
		},
		{
			name:      "too long",
			verifier:  string(make([]byte, 129)),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCodeVerifier(tt.verifier)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateCodeVerifier() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
