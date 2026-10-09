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
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func validateGetAuthOptionsByEmailRequest(
	req *authproto.GetAuthOptionsByEmailRequest,
) error {
	if req.Email == "" {
		return auth.StatusInvalidArguments.Err()
	}
	return nil
}

func validateCreateDomainAuthPolicyRequest(
	req *authproto.CreateDomainAuthPolicyRequest,
) error {
	if req.Domain == "" {
		return auth.StatusInvalidArguments.Err()
	}

	if req.AuthPolicy == nil {
		return auth.StatusInvalidArguments.Err()
	}

	return nil
}

func validateUpdateDomainAuthPolicyRequest(
	req *authproto.UpdateDomainAuthPolicyRequest,
) error {
	if req.Domain == "" {
		return auth.StatusInvalidArguments.Err()
	}

	if req.AuthPolicy == nil {
		return auth.StatusInvalidArguments.Err()
	}

	return nil
}

func validateGetDomainAuthPolicyRequest(
	req *authproto.GetDomainAuthPolicyRequest,
) error {
	if req.Domain == "" {
		return auth.StatusInvalidArguments.Err()
	}
	return nil
}

func validateDeleteDomainAuthPolicyRequest(
	req *authproto.DeleteDomainAuthPolicyRequest,
) error {
	if req.Domain == "" {
		return auth.StatusInvalidArguments.Err()
	}
	return nil
}

func validateListDomainAuthPoliciesRequest(
	req *authproto.ListDomainAuthPoliciesRequest,
) error {
	// All fields are optional for list request
	return nil
}
