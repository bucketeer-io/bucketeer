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

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Authenticator interface {
	Login(
		ctx context.Context,
		state, redirectURL string,
	) (string, error)
	Exchange(
		ctx context.Context,
		code, redirectURL string,
	) (*UserInfo, error)
}

type UserInfo struct {
	Name          string `json:"name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GoogleConfig struct {
	Issuer       string   `json:"issuer"`
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	RedirectURLs []string `json:"redirectUrls"`
}

type DemoSignInConfig struct {
	Enabled                bool   `json:"enabled"`
	Password               string `json:"password"`
	Email                  string `json:"email"`
	OrganizationId         string `json:"organizationId"`
	OrganizationOwnerEmail string `json:"organizationOwnerEmail"`
	ProjectId              string `json:"projectId"`
	EnvironmentId          string `json:"environmentId"`
	IsSystemAdmin          bool   `json:"isSystemAdmin"`
}

type PasswordPolicyConfig struct {
	MinLength        int  `json:"minLength"`
	RequireUppercase bool `json:"requireUppercase"`
	RequireLowercase bool `json:"requireLowercase"`
	RequireNumbers   bool `json:"requireNumbers"`
	RequireSymbols   bool `json:"requireSymbols"`
}

type PasswordTokensConfig struct {
	ResetTTL time.Duration `json:"resetTTL"`
	SetupTTL time.Duration `json:"setupTTL"`
}

type PasswordURLsConfig struct {
	ResetPath  string `json:"resetPath"`  // Path for password reset page
	SetupPath  string `json:"setupPath"`  // Path for password setup page
	TokenParam string `json:"tokenParam"` // URL parameter name for token
}

type PasswordAuthConfig struct {
	Enabled bool                 `json:"enabled"`
	Policy  PasswordPolicyConfig `json:"policy"`
	Tokens  PasswordTokensConfig `json:"tokens"`
	URLs    PasswordURLsConfig   `json:"urls"`
}

// UnmarshalJSON implements custom JSON unmarshaling for PasswordAuthConfig
// to handle duration strings like "1h", "24h" in the Tokens field
func (c *PasswordAuthConfig) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with string duration fields in Tokens
	type TokensAlias struct {
		ResetTTL string `json:"resetTTL"`
		SetupTTL string `json:"setupTTL"`
	}

	type Alias struct {
		Enabled bool                 `json:"enabled"`
		Policy  PasswordPolicyConfig `json:"policy"`
		Tokens  TokensAlias          `json:"tokens"`
		URLs    PasswordURLsConfig   `json:"urls"`
	}

	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set non-duration values
	c.Enabled = aux.Enabled
	c.Policy = aux.Policy
	c.URLs = aux.URLs

	// Only parse TTL values if password auth is enabled
	if aux.Enabled {
		resetTTL, err := time.ParseDuration(aux.Tokens.ResetTTL)
		if err != nil {
			return fmt.Errorf("failed to parse resetTTL: %w", err)
		}

		setupTTL, err := time.ParseDuration(aux.Tokens.SetupTTL)
		if err != nil {
			return fmt.Errorf("failed to parse setupTTL: %w", err)
		}

		c.Tokens.ResetTTL = resetTTL
		c.Tokens.SetupTTL = setupTTL
	}

	return nil
}

type OAuthConfig struct {
	Issuer       string             `json:"issuer"`
	Audience     string             `json:"audience"`
	GoogleConfig GoogleConfig       `json:"google"`
	Password     PasswordAuthConfig `json:"password"`
	DemoSignIn   DemoSignInConfig   `json:"demoSignIn"`
}
