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

package authz

import (
	"context"

	"google.golang.org/grpc/metadata"

	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// Metadata keys the public API gateway sets on the requests it relays.
const (
	APIKeyTokenMDKey      = "apikey-token"
	APIKeyMaintainerMDKey = "apikey-maintainer"
	APIKeyNameMDKey       = "apikey-name"
)

// Caller is who is asking, as the authorizer sees them. A request relayed
// by the public API gateway is attributed to the API key's maintainer, with
// the maintainer's own membership and roles.
type Caller struct {
	Email          string
	Name           string
	IsSystemAdmin  bool
	IsServiceToken bool
	// PublicAPI is set when the request was relayed by the public API gateway.
	PublicAPI *eventproto.Editor_PublicAPIEditor
}

// FullAccess reports whether the caller bypasses every check: a service
// token (machine account of batch and subscriber jobs) that is not relaying
// a public API request.
func (c *Caller) FullAccess() bool {
	return c.IsSystemAdmin && c.IsServiceToken && c.PublicAPI == nil
}

// callerFromContext reads the verified access token, and the API key
// metadata when a service token relays a public API request. API key
// metadata is trusted on service token requests only, because only the
// public API gateway sets it.
func callerFromContext(ctx context.Context) (*Caller, bool) {
	token, ok := rpc.GetAccessToken(ctx)
	if !ok {
		return nil, false
	}
	caller := &Caller{
		Email:          token.Email,
		Name:           token.Name,
		IsSystemAdmin:  token.IsSystemAdmin,
		IsServiceToken: token.IsServiceToken,
	}
	if !token.IsSystemAdmin || !token.IsServiceToken {
		return caller, true
	}
	if publicAPI := apiKeyEditorFromContext(ctx); publicAPI != nil && publicAPI.Token != "" {
		return &Caller{Email: publicAPI.Maintainer, PublicAPI: publicAPI}, true
	}
	return caller, true
}

func apiKeyEditorFromContext(ctx context.Context) *eventproto.Editor_PublicAPIEditor {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	apiKeyToken := md.Get(APIKeyTokenMDKey)
	if len(apiKeyToken) == 0 {
		return nil
	}
	editor := &eventproto.Editor_PublicAPIEditor{
		// Only the obfuscated key may reach the domain events and the audit log.
		Token: accdomain.ObfuscateAPIKey(apiKeyToken[0]),
	}
	if v := md.Get(APIKeyMaintainerMDKey); len(v) > 0 {
		editor.Maintainer = v[0]
	}
	if v := md.Get(APIKeyNameMDKey); len(v) > 0 {
		editor.Name = v[0]
	}
	return editor
}
