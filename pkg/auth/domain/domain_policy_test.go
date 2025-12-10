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
	"testing"

	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid email lowercase",
			email:   "user@example.com",
			want:    "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email uppercase",
			email:   "USER@EXAMPLE.COM",
			want:    "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with spaces",
			email:   "  user@example.com  ",
			want:    "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email mixed case",
			email:   "UsEr@ExAmPlE.CoM",
			want:    "user@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "no @ symbol",
			email:   "userexample.com",
			want:    "",
			wantErr: true,
		},
		{
			name:    "starts with @",
			email:   "@example.com",
			want:    "",
			wantErr: true,
		},
		{
			name:    "ends with @",
			email:   "user@",
			want:    "",
			wantErr: true,
		},
		{
			name:    "multiple @ symbols",
			email:   "user@@example.com",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "valid email uppercase",
			email:   "USER@EXAMPLE.COM",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			want:    "mail.example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "no @ symbol",
			email:   "userexample.com",
			want:    "",
			wantErr: true,
		},
		{
			name:    "domain with spaces",
			email:   "user@exam ple.com",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractDomain(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAuthPolicy_IsPasswordEnabled(t *testing.T) {
	tests := []struct {
		name   string
		policy *DomainAuthPolicy
		want   bool
	}{
		{
			name: "password enabled",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{
					Password: &authproto.PasswordAuthOption{
						Enabled: true,
					},
				},
			},
			want: true,
		},
		{
			name: "password disabled",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{
					Password: &authproto.PasswordAuthOption{
						Enabled: false,
					},
				},
			},
			want: false,
		},
		{
			name: "no password option",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{},
			},
			want: false,
		},
		{
			name:   "no auth policy",
			policy: &DomainAuthPolicy{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.IsPasswordEnabled(); got != tt.want {
				t.Errorf("DomainAuthPolicy.IsPasswordEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAuthPolicy_IsGoogleOidcEnabled(t *testing.T) {
	tests := []struct {
		name   string
		policy *DomainAuthPolicy
		want   bool
	}{
		{
			name: "google oidc enabled",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{
					GoogleOidc: &authproto.GoogleOidcOption{
						Enabled: true,
					},
				},
			},
			want: true,
		},
		{
			name: "google oidc disabled",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{
					GoogleOidc: &authproto.GoogleOidcOption{
						Enabled: false,
					},
				},
			},
			want: false,
		},
		{
			name: "no google oidc option",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{},
			},
			want: false,
		},
		{
			name:   "no auth policy",
			policy: &DomainAuthPolicy{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.IsGoogleOidcEnabled(); got != tt.want {
				t.Errorf("DomainAuthPolicy.IsGoogleOidcEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAuthPolicy_IsCompanyOidcRequired(t *testing.T) {
	tests := []struct {
		name   string
		policy *DomainAuthPolicy
		want   bool
	}{
		{
			name: "company oidc required",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{
					CompanyOidc: &authproto.CompanyOidcOption{
						Enabled:  true,
						Required: true,
					},
				},
			},
			want: true,
		},
		{
			name: "company oidc not required",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{
					CompanyOidc: &authproto.CompanyOidcOption{
						Enabled:  true,
						Required: false,
					},
				},
			},
			want: false,
		},
		{
			name: "no company oidc option",
			policy: &DomainAuthPolicy{
				AuthPolicy: &authproto.AuthPolicy{},
			},
			want: false,
		},
		{
			name:   "no auth policy",
			policy: &DomainAuthPolicy{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.IsCompanyOidcRequired(); got != tt.want {
				t.Errorf("DomainAuthPolicy.IsCompanyOidcRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAuthPolicy_ToProto(t *testing.T) {
	policy := &DomainAuthPolicy{
		Domain: "example.com",
		AuthPolicy: &authproto.AuthPolicy{
			Password: &authproto.PasswordAuthOption{
				Enabled: true,
			},
		},
		Enabled:   true,
		CreatedAt: 123456789,
		UpdatedAt: 123456789,
	}

	proto := policy.ToProto()

	if proto.Domain != policy.Domain {
		t.Errorf("ToProto().Domain = %v, want %v", proto.Domain, policy.Domain)
	}
	if proto.Enabled != policy.Enabled {
		t.Errorf("ToProto().Enabled = %v, want %v", proto.Enabled, policy.Enabled)
	}
	if proto.CreatedAt != policy.CreatedAt {
		t.Errorf("ToProto().CreatedAt = %v, want %v", proto.CreatedAt, policy.CreatedAt)
	}
	if proto.UpdatedAt != policy.UpdatedAt {
		t.Errorf("ToProto().UpdatedAt = %v, want %v", proto.UpdatedAt, policy.UpdatedAt)
	}
}

func TestNewDomainAuthPolicyFromProto(t *testing.T) {
	proto := &authproto.DomainAuthPolicy{
		Domain: "example.com",
		AuthPolicy: &authproto.AuthPolicy{
			Password: &authproto.PasswordAuthOption{
				Enabled: true,
			},
		},
		Enabled:   true,
		CreatedAt: 123456789,
		UpdatedAt: 123456789,
	}

	policy := NewDomainAuthPolicyFromProto(proto)

	if policy.Domain != proto.Domain {
		t.Errorf("NewDomainAuthPolicyFromProto().Domain = %v, want %v", policy.Domain, proto.Domain)
	}
	if policy.Enabled != proto.Enabled {
		t.Errorf("NewDomainAuthPolicyFromProto().Enabled = %v, want %v", policy.Enabled, proto.Enabled)
	}
	if policy.CreatedAt != proto.CreatedAt {
		t.Errorf("NewDomainAuthPolicyFromProto().CreatedAt = %v, want %v", policy.CreatedAt, proto.CreatedAt)
	}
	if policy.UpdatedAt != proto.UpdatedAt {
		t.Errorf("NewDomainAuthPolicyFromProto().UpdatedAt = %v, want %v", policy.UpdatedAt, proto.UpdatedAt)
	}

	// Test nil proto
	nilPolicy := NewDomainAuthPolicyFromProto(nil)
	if nilPolicy != nil {
		t.Errorf("NewDomainAuthPolicyFromProto(nil) = %v, want nil", nilPolicy)
	}
}
