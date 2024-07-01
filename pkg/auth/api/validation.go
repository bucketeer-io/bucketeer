//  Copyright 2024 The Bucketeer Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package api

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func validateGetAuthenticationURLRequest(
	req *authproto.GetAuthenticationURLRequest,
	localizer locale.Localizer,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		dt, err := auth.StatusMissingAuthType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "auth_type"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.State == "" {
		dt, err := auth.StatusMissingState.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "state"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := auth.StatusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateExchangeBucketeerTokenRequest(
	req *authproto.ExchangeBucketeerTokenRequest,
	localizer locale.Localizer,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		dt, err := auth.StatusMissingAuthType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "auth_type"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.Code == "" {
		dt, err := auth.StatusMissingCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "code"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := auth.StatusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateRefreshBucketeerTokenRequest(
	req *authproto.RefreshBucketeerTokenRequest,
	localizer locale.Localizer,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		dt, err := auth.StatusMissingAuthType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "auth_type"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.RefreshToken == "" {
		dt, err := auth.StatusMissingRefreshToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "refresh_token"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := auth.StatusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}
