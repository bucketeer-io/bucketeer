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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestVerifyEmailFormat(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		email string
		ok    bool
	}{
		{"foo@gmail.com", true},
		{"foo+bar@abc.co.jp", true},
		{"invalid", false},
		{"@invalid", false},
		{"", false},
	}
	for _, tc := range testcases {
		ok := verifyEmailFormat(tc.email)
		assert.Equal(t, tc.ok, ok, tc.email)
	}
}

func TestValidateCreateAccountV2Request(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		req         *accountproto.CreateAccountV2Request
		expectedErr error
	}{
		{
			desc:        "err: missing organization id",
			req:         &accountproto.CreateAccountV2Request{},
			expectedErr: statusMissingOrganizationID.Err(),
		},
		{
			desc: "err: missing email",
			req: &accountproto.CreateAccountV2Request{
				OrganizationId: "org-id",
			},
			expectedErr: statusEmailIsEmpty.Err(),
		},
		{
			desc: "err: invalid email",
			req: &accountproto.CreateAccountV2Request{
				OrganizationId: "org-id",
				Email:          "invalid-email",
			},
			expectedErr: statusInvalidEmail.Err(),
		},
		{
			desc: "err: invalid organization role",
			req: &accountproto.CreateAccountV2Request{
				OrganizationId:   "org-id",
				Email:            "test@example.com",
				OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
			},
			expectedErr: statusInvalidOrganizationRole.Err(),
		},
		{
			desc: "err: missing environment roles for member",
			req: &accountproto.CreateAccountV2Request{
				OrganizationId:   "org-id",
				Email:            "test@example.com",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			},
			expectedErr: statusInvalidEnvironmentRole.Err(),
		},
		{
			desc: "success: admin role without environment roles",
			req: &accountproto.CreateAccountV2Request{
				OrganizationId:   "org-id",
				Email:            "test@example.com",
				OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			},
			expectedErr: nil,
		},
		{
			desc: "success: member role with environment roles",
			req: &accountproto.CreateAccountV2Request{
				OrganizationId:   "org-id",
				Email:            "test@example.com",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env-id", Role: accountproto.AccountV2_Role_Environment_VIEWER},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateCreateAccountV2Request(p.req)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDeleteAccountV2Request(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		req         *accountproto.DeleteAccountV2Request
		expectedErr error
	}{
		{
			desc:        "err: missing email",
			req:         &accountproto.DeleteAccountV2Request{OrganizationId: "org-id"},
			expectedErr: statusEmailIsEmpty.Err(),
		},
		{
			desc:        "err: invalid email",
			req:         &accountproto.DeleteAccountV2Request{Email: "invalid", OrganizationId: "org-id"},
			expectedErr: statusInvalidEmail.Err(),
		},
		{
			desc:        "err: missing organization id",
			req:         &accountproto.DeleteAccountV2Request{Email: "test@example.com"},
			expectedErr: statusMissingOrganizationID.Err(),
		},
		{
			desc:        "success",
			req:         &accountproto.DeleteAccountV2Request{Email: "test@example.com", OrganizationId: "org-id"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateDeleteAccountV2Request(p.req)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateAccountV2Request(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		req      *accountproto.UpdateAccountV2Request
		expected error
	}{
		{
			desc:     "err: missing email",
			req:      &accountproto.UpdateAccountV2Request{OrganizationId: "org-id"},
			expected: statusEmailIsEmpty.Err(),
		},
		{
			desc:     "success",
			req:      &accountproto.UpdateAccountV2Request{Email: "email@example.com", OrganizationId: "org-id"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateUpdateAccountV2Request(p.req)
			if p.expected != nil {
				assert.Equal(t, p.expected.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEnableAccountV2Request(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		req      *accountproto.EnableAccountV2Request
		expected error
	}{
		{
			desc:     "err: missing email",
			req:      &accountproto.EnableAccountV2Request{OrganizationId: "org-id"},
			expected: statusEmailIsEmpty.Err(),
		},
		{
			desc:     "success",
			req:      &accountproto.EnableAccountV2Request{Email: "email@example.com", OrganizationId: "org-id"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateEnableAccountV2Request(p.req)
			if p.expected != nil {
				assert.Equal(t, p.expected.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDisableAccountV2Request(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		req      *accountproto.DisableAccountV2Request
		expected error
	}{
		{
			desc:     "err: missing email",
			req:      &accountproto.DisableAccountV2Request{OrganizationId: "org-id"},
			expected: statusEmailIsEmpty.Err(),
		},
		{
			desc:     "success",
			req:      &accountproto.DisableAccountV2Request{Email: "email@example.com", OrganizationId: "org-id"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateDisableAccountV2Request(p.req)
			if p.expected != nil {
				assert.Equal(t, p.expected.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGetAccountV2Request(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		req      *accountproto.GetAccountV2Request
		expected error
	}{
		{
			desc:     "err: missing email",
			req:      &accountproto.GetAccountV2Request{OrganizationId: "org-id"},
			expected: statusEmailIsEmpty.Err(),
		},
		{
			desc:     "success",
			req:      &accountproto.GetAccountV2Request{Email: "email@example.com", OrganizationId: "org-id"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateGetAccountV2Request(p.req)
			if p.expected != nil {
				assert.Equal(t, p.expected.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGetAccountV2ByEnvironmentIDRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		req      *accountproto.GetAccountV2ByEnvironmentIDRequest
		expected error
	}{
		{
			desc:     "err: missing email",
			req:      &accountproto.GetAccountV2ByEnvironmentIDRequest{},
			expected: statusEmailIsEmpty.Err(),
		},
		{
			desc:     "success",
			req:      &accountproto.GetAccountV2ByEnvironmentIDRequest{Email: "email@example.com"},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateGetAccountV2ByEnvironmentIDRequest(p.req)
			if p.expected != nil {
				assert.Equal(t, p.expected.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
