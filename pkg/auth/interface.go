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

type EmailServiceConfig struct {
	Provider       string `json:"provider"` // "smtp", "sendgrid", "ses"
	SMTPHost       string `json:"smtpHost"`
	SMTPPort       int    `json:"smtpPort"`
	SMTPUsername   string `json:"smtpUsername"`
	SMTPPassword   string `json:"smtpPassword"`
	SendGridAPIKey string `json:"sendgridAPIKey"`
	SESRegion      string `json:"sesRegion"`
	SESAccessKey   string `json:"sesAccessKey"`
	SESSecretKey   string `json:"sesSecretKey"`
	FromEmail      string `json:"fromEmail"`
	FromName       string `json:"fromName"`
	BaseURL        string `json:"baseURL"` // For constructing reset URLs
}

type PasswordAuthConfig struct {
	Enabled                  bool               `json:"enabled"`
	PasswordMinLength        int                `json:"passwordMinLength"`
	PasswordRequireUppercase bool               `json:"passwordRequireUppercase"`
	PasswordRequireLowercase bool               `json:"passwordRequireLowercase"`
	PasswordRequireNumbers   bool               `json:"passwordRequireNumbers"`
	PasswordRequireSymbols   bool               `json:"passwordRequireSymbols"`
	PasswordResetTokenTTL    time.Duration      `json:"passwordResetTokenTTL"`
	EmailServiceEnabled      bool               `json:"emailServiceEnabled"`
	EmailServiceConfig       EmailServiceConfig `json:"emailServiceConfig"`
}

type OAuthConfig struct {
	Issuer       string             `json:"issuer"`
	Audience     string             `json:"audience"`
	GoogleConfig GoogleConfig       `json:"google"`
	DemoSignIn   DemoSignInConfig   `json:"demoSignIn"`
	PasswordAuth PasswordAuthConfig `json:"passwordAuth"`
}
