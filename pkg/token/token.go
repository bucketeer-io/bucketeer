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

package token

import (
	"time"
)

type AccessToken struct {
	Issuer         string    `json:"iss"`
	Audience       string    `json:"aud"`
	Expiry         time.Time `json:"exp"`
	IssuedAt       time.Time `json:"iat"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	IsSystemAdmin  bool      `json:"is_system_admin"`
	OrganizationID string    `json:"organization_id"`
}

type RefreshToken struct {
	Email          string    `json:"email"`
	Expiry         time.Time `json:"exp"`
	IssuedAt       time.Time `json:"iat"`
	OrganizationID string    `json:"organization_id"`
}

type DemoCreationToken struct {
	Issuer   string    `json:"iss"`
	Audience string    `json:"aud"`
	Expiry   time.Time `json:"exp"`
	IssuedAt time.Time `json:"iat"`
	Email    string    `json:"email"`
}
