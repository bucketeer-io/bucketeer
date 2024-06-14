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

package google

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	defaultScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
)

type Authenticator struct {
	config        *oauth2.Config
	accountClient accountclient.Client
	signer        token.Signer
	emailFilter   *regexp.Regexp
	logger        *zap.Logger
}

func NewAuthenticator(
	config auth.GoogleConfig,
	accountClient accountclient.Client,
	signer token.Signer,
	logger *zap.Logger,
) *Authenticator {
	return &Authenticator{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes:       defaultScopes,
			Endpoint:     google.Endpoint,
		},
		accountClient: accountClient,
		signer:        signer,
		logger:        logger.Named("google-authenticator"),
	}
}

func (a Authenticator) Login(ctx context.Context, state string, localizer locale.Localizer) string {
	return a.config.AuthCodeURL(state)
}

func (a Authenticator) Exchange(
	ctx context.Context,
	code string,
	localizer locale.Localizer,
) (*authproto.Token, error) {
	authToken, err := a.config.Exchange(ctx, code)
	if err != nil {
		a.logger.Error("Google: failed to exchange token", zap.Error(err))
		return nil, err
	}
	return a.generateToken(ctx, authToken, localizer)
}

func (a Authenticator) Refresh(
	ctx context.Context,
	token string,
	expires time.Duration,
	localizer locale.Localizer,
) (*authproto.Token, error) {
	t := &oauth2.Token{
		RefreshToken: token,
		Expiry:       time.Now().Add(expires),
	}
	newToken, err := a.config.TokenSource(ctx, t).Token()
	if err != nil {
		a.logger.Error("Google: failed to refresh token", zap.Error(err))
		return nil, err
	}
	return a.generateToken(ctx, newToken, localizer)
}

func (a Authenticator) generateToken(
	ctx context.Context,
	t *oauth2.Token,
	localizer locale.Localizer,
) (*authproto.Token, error) {
	client := a.config.Client(ctx, t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		a.logger.Error("Google: failed to get user info", zap.Error(err))
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			a.logger.Error("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	var userInfo struct {
		Username      string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		a.logger.Error("Failed to decode user info", zap.Error(err))
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}

	if err := a.maybeCheckEmail(ctx, userInfo.Email, localizer); err != nil {
		a.logger.Info(
			"Access denied email",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", userInfo.Email))...,
		)
		return nil, err
	}

	orgResp, err := a.accountClient.GetMyOrganizationsByEmail(
		ctx,
		&accountproto.GetMyOrganizationsByEmailRequest{
			Email: userInfo.Email,
		},
	)
	if err != nil {
		a.logger.Error(
			"Failed to get account's organizations",
			zap.Error(err),
			zap.String("email", userInfo.Email),
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	if len(orgResp.Organizations) == 0 {
		a.logger.Error(
			"Unable to generate token for an unapproved account",
			zap.String("email", userInfo.Email),
		)
		dt, err := auth.StatusUnapprovedAccount.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}

	idToken := &token.IDToken{
		Expiry:        t.Expiry,
		Email:         userInfo.Email,
		IsSystemAdmin: hasSystemAdminOrganization(orgResp.Organizations),
	}
	signedIDToken, err := a.signer.Sign(idToken)
	if err != nil {
		a.logger.Error(
			"Failed to sign id token",
			zap.Error(err),
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}

	return &authproto.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry.Unix(),
		IdToken:      signedIDToken,
	}, nil
}

func (a Authenticator) maybeCheckEmail(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) error {
	if a.emailFilter == nil {
		return nil
	}
	if a.emailFilter.MatchString(email) {
		return nil
	}
	dt, err := auth.StatusAccessDeniedEmail.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.PermissionDenied),
	})
	if err != nil {
		return auth.StatusInternal.Err()
	}
	return dt.Err()
}

func hasSystemAdminOrganization(orgs []*environmentproto.Organization) bool {
	for _, org := range orgs {
		if org.SystemAdmin {
			return true
		}
	}
	return false
}
