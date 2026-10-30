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
)

type Authenticator interface {
	Login(
		ctx context.Context,
		state, redirectURL string,
	) (string, error)
	// Exchange trades the authorization code for the identity of the signed-in
	// user. The returned email address is used as the account key, so
	// implementations must only return identities whose email address the
	// identity provider has verified.
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
	// AllowedDomains restricts sign-in to Google Workspace accounts whose
	// hosted domain (the "hd" claim reported by Google) is in this list.
	// When empty, any verified Google account may sign in.
	AllowedDomains []string `json:"allowedDomains"`
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

type OAuthConfig struct {
	Issuer       string           `json:"issuer"`
	Audience     string           `json:"audience"`
	GoogleConfig GoogleConfig     `json:"google"`
	DemoSignIn   DemoSignInConfig `json:"demoSignIn"`
}
