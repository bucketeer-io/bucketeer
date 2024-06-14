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

package auth

import (
	"context"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

type Authenticator interface {
	Login(ctx context.Context, state string, localizer locale.Localizer) string
	Exchange(ctx context.Context, code string, localizer locale.Localizer) (*authproto.Token, error)
	Refresh(
		ctx context.Context,
		token string,
		expires time.Duration,
		localizer locale.Localizer,
	) (*authproto.Token, error)
}

type GoogleConfig struct {
	Issuer       string `json:"issuer"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectUrl"`
}

type OAuthConfig struct {
	GoogleConfig GoogleConfig `json:"google"`
}
