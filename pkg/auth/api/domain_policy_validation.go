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
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func validateGetAuthOptionsByEmailRequest(
	req *authproto.GetAuthOptionsByEmailRequest,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}
	return nil
}

func validateCreateDomainAuthPolicyRequest(
	req *authproto.CreateDomainAuthPolicyRequest,
	localizer locale.Localizer,
) error {
	if req.Domain == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "domain"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.AuthPolicy == nil {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "auth_policy"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	return nil
}

func validateUpdateDomainAuthPolicyRequest(
	req *authproto.UpdateDomainAuthPolicyRequest,
	localizer locale.Localizer,
) error {
	if req.Domain == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "domain"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.AuthPolicy == nil {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "auth_policy"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	return nil
}

func validateGetDomainAuthPolicyRequest(
	req *authproto.GetDomainAuthPolicyRequest,
	localizer locale.Localizer,
) error {
	if req.Domain == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "domain"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}
	return nil
}

func validateDeleteDomainAuthPolicyRequest(
	req *authproto.DeleteDomainAuthPolicyRequest,
	localizer locale.Localizer,
) error {
	if req.Domain == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "domain"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}
	return nil
}

func validateListDomainAuthPoliciesRequest(
	req *authproto.ListDomainAuthPoliciesRequest,
	localizer locale.Localizer,
) error {
	// All fields are optional for list request
	return nil
}
