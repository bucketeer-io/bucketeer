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

package domain

import (
	"time"
)

// EmailVerificationToken represents an email verification token for magic link authentication
type EmailVerificationToken struct {
	Email      string
	Token      string
	CreatedAt  int64
	ExpiresAt  int64
	VerifiedAt *int64 // Nullable: NULL = not verified, value = verification timestamp
	IPAddress  string
	UserAgent  string
}

// IsExpired checks if the email verification token has expired
func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().Unix() > t.ExpiresAt
}

// IsVerified checks if the token has already been verified
func (t *EmailVerificationToken) IsVerified() bool {
	return t.VerifiedAt != nil
}

// IsValid checks if the token is valid (not expired and not yet verified)
func (t *EmailVerificationToken) IsValid() bool {
	return t.Token != "" && !t.IsExpired() && !t.IsVerified()
}

// WasRecentlyVerified checks if token was verified within the last 5 minutes
// This allows users who click the link multiple times to still see org selection
func (t *EmailVerificationToken) WasRecentlyVerified() bool {
	if t.VerifiedAt == nil {
		return false
	}
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute).Unix()
	return *t.VerifiedAt > fiveMinutesAgo
}
