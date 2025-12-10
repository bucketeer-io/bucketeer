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

package rpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
)

type contextKey int

const (
	AccessTokenKey contextKey = iota
	DemoCreationTokenKey
)

const (
	healthServiceName          = "/grpc.health.v1.Health/"
	flagTriggerWebhookName     = "/bucketeer.feature.FeatureService/FlagTriggerWebhook"
	exchangeDemoTokenName      = "/bucketeer.environment.EnvironmentService/ExchangeDemoToken"
	createDemoOrganizationName = "/bucketeer.environment.EnvironmentService/CreateDemoOrganization"
	initiatePasswordResetName  = "/bucketeer.auth.AuthService/InitiatePasswordReset"
	resetPasswordName          = "/bucketeer.auth.AuthService/ResetPassword"
	validatePasswordResetName  = "/bucketeer.auth.AuthService/ValidatePasswordResetToken"
	initiatePasswordSetupName  = "/bucketeer.auth.AuthService/InitiatePasswordSetup"
	setupPasswordName          = "/bucketeer.auth.AuthService/SetupPassword"
	validatePasswordSetupName  = "/bucketeer.auth.AuthService/ValidatePasswordSetupToken"
	// New login methods
	signInPasswordName           = "/bucketeer.auth.AuthService/SignInPassword"
	switchOrganizationName       = "/bucketeer.auth.AuthService/SwitchOrganization"
	getAuthOptionsByEmailName    = "/bucketeer.auth.AuthService/GetAuthOptionsByEmail"
	getGoogleOidcAuthURLName     = "/bucketeer.auth.AuthService/GetGoogleOidcAuthURL"
	exchangeGoogleOidcTokenName  = "/bucketeer.auth.AuthService/ExchangeGoogleOidcToken"
	getCompanyOidcAuthURLName    = "/bucketeer.auth.AuthService/GetCompanyOidcAuthURL"
	exchangeCompanyOidcTokenName = "/bucketeer.auth.AuthService/ExchangeCompanyOidcToken"
	getDemoSiteStatusName        = "/bucketeer.auth.AuthService/GetDemoSiteStatus"
	// Old login methods (for backward compatibility)
	signInName               = "/bucketeer.auth.AuthService/SignIn"
	getAuthenticationURLName = "/bucketeer.auth.AuthService/GetAuthenticationURL"
	exchangeTokenName        = "/bucketeer.auth.AuthService/ExchangeToken"
	refreshTokenName         = "/bucketeer.auth.AuthService/RefreshToken"
)

type authFunc func(verifier token.Verifier, token string) (interface{}, error)

type methodAuth struct {
	authFunc authFunc
	key      interface{}
}

var specificAuthMethods = map[string]methodAuth{
	createDemoOrganizationName: {
		authFunc: func(v token.Verifier, token string) (interface{}, error) {
			return v.VerifyDemoCreationToken(token)
		},
		key: DemoCreationTokenKey,
	},
}

var defaultAuth = methodAuth{
	authFunc: func(v token.Verifier, token string) (interface{}, error) {
		return v.VerifyAccessToken(token)
	},
	key: AccessTokenKey,
}

var (
	skipAuthMethods = []string{
		healthServiceName,
		flagTriggerWebhookName,
		exchangeDemoTokenName,
		initiatePasswordResetName,
		resetPasswordName,
		validatePasswordResetName,
		initiatePasswordSetupName,
		setupPasswordName,
		validatePasswordSetupName,
		// New login methods
		signInPasswordName,
		switchOrganizationName,
		getAuthOptionsByEmailName,
		getGoogleOidcAuthURLName,
		exchangeGoogleOidcTokenName,
		getCompanyOidcAuthURLName,
		exchangeCompanyOidcTokenName,
		getDemoSiteStatusName,
		// Old login methods (for backward compatibility)
		signInName,
		getAuthenticationURLName,
		exchangeTokenName,
		refreshTokenName,
	}
)

func AuthUnaryServerInterceptor(verifier token.Verifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		for _, method := range skipAuthMethods {
			if strings.HasPrefix(info.FullMethod, method) {
				return handler(ctx, req)
			}
		}
		authConfig := defaultAuth
		for method, config := range specificAuthMethods {
			if strings.HasPrefix(info.FullMethod, method) {
				authConfig = config
				break
			}
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "token is required")
		}
		rawTokens, ok := md["authorization"]
		if !ok || len(rawTokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token is required")
		}
		subs := strings.Split(rawTokens[0], " ")
		if len(subs) != 2 {
			return nil, status.Error(codes.Unauthenticated, "token is malformed")
		}
		token, err := authConfig.authFunc(verifier, subs[1])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %s", err.Error())
		}
		ctx = context.WithValue(ctx, authConfig.key, token)
		return handler(ctx, req)
	}
}

func GetAccessToken(ctx context.Context) (*token.AccessToken, bool) {
	t, ok := ctx.Value(AccessTokenKey).(*token.AccessToken)
	return t, ok
}

func GetDemoCreationToken(ctx context.Context) (*token.DemoCreationToken, bool) {
	t, ok := ctx.Value(DemoCreationTokenKey).(*token.DemoCreationToken)
	return t, ok
}
