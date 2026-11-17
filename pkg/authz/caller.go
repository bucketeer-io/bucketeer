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
	OrganizationID string
	// Service is set for the machine identity of batch and subscriber jobs.
	Service bool
	// PublicAPI is set when the request was relayed by the public API gateway.
	PublicAPI *eventproto.Editor_PublicAPIEditor
}

// FullAccess reports whether the caller skips the permission checks: a
// service identity that is not relaying a public API request.
func (c Caller) FullAccess() bool {
	return c.Service && c.PublicAPI == nil
}

// Relayed returns the API key's maintainer as the caller when the public API
// gateway relayed the request. Only the gateway relays with a service
// identity, so the metadata is read on service callers only.
func (c Caller) Relayed(ctx context.Context) Caller {
	if !c.Service {
		return c
	}
	publicAPI := apiKeyEditorFromContext(ctx)
	if publicAPI == nil {
		return c
	}
	return Caller{Email: publicAPI.Maintainer, PublicAPI: publicAPI}
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
