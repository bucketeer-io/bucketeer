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

package domain

import (
	"time"
)

// AccountCredentials represents password credentials for an account
type AccountCredentials struct {
	Email        string
	PasswordHash string
	CreatedAt    int64
	UpdatedAt    int64
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	Token     string
	Email     string
	ExpiresAt int64
	CreatedAt int64
}

// IsExpired checks if the password reset token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().Unix() > t.ExpiresAt
}

// IsValid checks if the password reset token is valid (not expired and not empty)
func (t *PasswordResetToken) IsValid() bool {
	return t.Token != "" && !t.IsExpired()
}
