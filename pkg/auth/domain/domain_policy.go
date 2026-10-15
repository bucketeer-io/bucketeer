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
	"errors"
	"strings"

	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

var (
	ErrInvalidEmail  = errors.New("domain: invalid email address")
	ErrInvalidDomain = errors.New("domain: invalid domain")
)

// DomainAuthPolicy represents a domain-based authentication policy
type DomainAuthPolicy struct {
	Domain     string
	AuthPolicy *authproto.AuthPolicy
	Enabled    bool
	CreatedAt  int64
	UpdatedAt  int64
}

// NormalizeEmail normalizes an email address for consistent processing
// Rules:
// - Convert to lowercase (email domains are case-insensitive per RFC 5321)
// - Trim leading/trailing whitespace
// - Reject empty strings
// - Basic format validation (must contain @)
//
// Note: We do NOT handle plus addressing (user+tag@domain.com) normalization
// as different domains may have different policies about this
func NormalizeEmail(email string) (string, error) {
	// Trim whitespace
	email = strings.TrimSpace(email)

	// Check if empty
	if email == "" {
		return "", ErrInvalidEmail
	}

	// Basic validation: must contain @
	if !strings.Contains(email, "@") {
		return "", ErrInvalidEmail
	}

	// Convert to lowercase (email addresses are case-insensitive)
	email = strings.ToLower(email)

	// Additional validation: must not start or end with @
	if strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return "", ErrInvalidEmail
	}

	// Validate that there's content before and after @
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", ErrInvalidEmail
	}

	return email, nil
}

// ExtractDomain extracts the domain from an email address
// The email should be normalized before calling this function
func ExtractDomain(email string) (string, error) {
	// Normalize first to ensure consistency
	normalizedEmail, err := NormalizeEmail(email)
	if err != nil {
		return "", err
	}

	// Split on @ and take the domain part
	parts := strings.Split(normalizedEmail, "@")
	if len(parts) != 2 {
		return "", ErrInvalidEmail
	}

	domain := parts[1]

	// Validate domain is not empty
	if domain == "" {
		return "", ErrInvalidDomain
	}

	// Basic domain validation: should not contain spaces or special chars except dots and hyphens
	if strings.Contains(domain, " ") {
		return "", ErrInvalidDomain
	}

	return domain, nil
}

// IsPasswordEnabled checks if password authentication is enabled for this policy
func (p *DomainAuthPolicy) IsPasswordEnabled() bool {
	if p.AuthPolicy == nil || p.AuthPolicy.Password == nil {
		return false
	}
	return p.AuthPolicy.Password.Enabled
}

// IsGoogleOidcEnabled checks if Google OIDC is enabled for this policy
func (p *DomainAuthPolicy) IsGoogleOidcEnabled() bool {
	if p.AuthPolicy == nil || p.AuthPolicy.GoogleOidc == nil {
		return false
	}
	return p.AuthPolicy.GoogleOidc.Enabled
}

// IsCompanyOidcEnabled checks if company OIDC is enabled for this policy
func (p *DomainAuthPolicy) IsCompanyOidcEnabled() bool {
	if p.AuthPolicy == nil || p.AuthPolicy.CompanyOidc == nil {
		return false
	}
	return p.AuthPolicy.CompanyOidc.Enabled
}

// IsCompanyOidcRequired checks if company OIDC is required (exclusive) for this policy
func (p *DomainAuthPolicy) IsCompanyOidcRequired() bool {
	if p.AuthPolicy == nil || p.AuthPolicy.CompanyOidc == nil {
		return false
	}
	return p.AuthPolicy.CompanyOidc.Required
}

// ToProto converts the domain entity to protobuf message
func (p *DomainAuthPolicy) ToProto() *authproto.DomainAuthPolicy {
	return &authproto.DomainAuthPolicy{
		Domain:     p.Domain,
		AuthPolicy: p.AuthPolicy,
		Enabled:    p.Enabled,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// NewDomainAuthPolicyFromProto creates a domain entity from protobuf message
func NewDomainAuthPolicyFromProto(proto *authproto.DomainAuthPolicy) *DomainAuthPolicy {
	if proto == nil {
		return nil
	}
	return &DomainAuthPolicy{
		Domain:     proto.Domain,
		AuthPolicy: proto.AuthPolicy,
		Enabled:    proto.Enabled,
		CreatedAt:  proto.CreatedAt,
		UpdatedAt:  proto.UpdatedAt,
	}
}
