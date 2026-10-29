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
	"encoding/json"
	"slices"
	"time"

	"github.com/jinzhu/copier"
	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var (
	ErrNameRequired = err.NewErrorInvalidArgEmpty(
		err.PushPackageName,
		"name is required",
		"name",
	)
	ErrTagRequired      = err.NewErrorInvalidArgEmpty(err.PushPackageName, "tag is required", "tag")
	ErrTagDuplicated    = err.NewErrorInvalidArgDuplicated(err.PushPackageName, "tag is duplicated", "tag")
	ErrTagAlreadyExists = err.NewErrorAlreadyExists(err.PushPackageName, "tag already exists")
	ErrTagNotFound      = err.NewErrorNotFound(err.PushPackageName, "tag not found", "tag")
)

type Push struct {
	*proto.Push
}

func NewPush(name, fcmServiceAccount string, tags []string) (*Push, error) {
	_, err := convMap(tags)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	p := &Push{&proto.Push{
		Name:              name,
		Id:                id.String(),
		FcmServiceAccount: fcmServiceAccount,
		Tags:              tags,
		Disabled:          false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}}
	return p, nil
}

func (p *Push) Update(
	name *wrapperspb.StringValue,
	tagChanges []*proto.TagChange,
	disabled *wrapperspb.BoolValue,
) (*Push, error) {
	updated := &Push{}
	if err := copier.Copy(updated, p); err != nil {
		return nil, err
	}
	if name != nil {
		updated.Name = name.Value
	}
	for _, change := range tagChanges {
		switch change.ChangeType {
		case proto.ChangeType_CREATE, proto.ChangeType_UPDATE:
			if err := updated.AddTag(change.Tag); err != nil {
				return nil, err
			}
		case proto.ChangeType_DELETE:
			if err := updated.RemoveTag(change.Tag); err != nil {
				return nil, err
			}
		}
	}
	if disabled != nil {
		updated.Disabled = disabled.Value
	}

	updated.UpdatedAt = time.Now().Unix()
	if err := validate(updated.Push); err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *Push) Rename(name string) error {
	p.Name = name
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// WithoutFCMServiceAccount returns a deep copy of the push with an empty FCM
// service account. The service account holds the private key used to send the
// notifications, so it is read only by the sender. Every other record of the
// push, such as the domain event and the audit log built from it, must use
// this copy so the key is never stored outside the push table.
func (p *Push) WithoutFCMServiceAccount() *Push {
	copied := &Push{Push: pb.Clone(p.Push).(*proto.Push)}
	copied.FcmServiceAccount = ""
	return copied
}

func convMap(ss []string) (map[string]int, error) {
	m := make(map[string]int, len(ss))
	for idx, t := range ss {
		if _, ok := m[t]; ok {
			return nil, ErrTagDuplicated
		}
		m[t] = idx
	}
	return m, nil
}

func (p *Push) AddTags(newTags []string) error {
	mtag, err := convMap(p.Tags)
	if err != nil {
		return err
	}
	for _, t := range newTags {
		if _, ok := mtag[t]; ok {
			return ErrTagAlreadyExists
		}
		p.Tags = append(p.Tags, t)
	}
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Push) DeleteTags(tags []string) error {
	tagMap, err := convMap(tags)
	if err != nil {
		return err
	}
	existMap, err := convMap(p.Tags)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if _, ok := existMap[t]; !ok {
			return ErrTagNotFound
		}
	}
	newTags := []string{}
	for _, t := range p.Tags {
		if _, ok := tagMap[t]; !ok {
			newTags = append(newTags, t)
		}
	}
	p.Tags = newTags
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Push) SetDeleted() {
	p.Deleted = true
	p.UpdatedAt = time.Now().Unix()
}

func (p *Push) ExistTag(findTag string) bool {
	for _, t := range p.Tags {
		if t == findTag {
			return true
		}
	}
	return false
}

func (p *Push) AddTag(tag string) error {
	if slices.Contains(p.Tags, tag) {
		return nil
	}
	p.Tags = append(p.Tags, tag)
	return nil
}

func (p *Push) RemoveTag(tag string) error {
	index := slices.Index(p.Tags, tag)
	if index == -1 {
		return ErrTagNotFound
	}
	p.Tags = slices.Delete(p.Tags, index, index+1)
	return nil
}

func validate(p *proto.Push) error {
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Tags == nil {
		return ErrTagRequired
	}
	return nil
}

// allowedFCMTokenURIs is the set of token endpoints accepted in a Firebase
// service-account key file. A valid key either omits token_uri or points at
// Google's OAuth token endpoint.
var allowedFCMTokenURIs = map[string]struct{}{
	"":                                    {},
	"https://oauth2.googleapis.com/token": {},
	"https://accounts.google.com/o/oauth2/token": {},
}

var (
	// ErrFCMServiceAccountRequired is returned when the FCM service account is empty.
	ErrFCMServiceAccountRequired = err.NewErrorInvalidArgEmpty(
		err.PushPackageName,
		"fcm service account is required",
		"fcmServiceAccount",
	)
	// ErrFCMServiceAccountInvalid is returned when the FCM service account is not a
	// Google service-account key file. Rejecting every other credential type is a
	// security control: google.CredentialsFromJSON honours the token_uri and, for
	// the "external_account" family, an attacker-controlled credential_source (a
	// local file or an arbitrary URL). A push editor could otherwise make the
	// subscriber read local files (e.g. the service token) or perform SSRF when the
	// token is exchanged.
	ErrFCMServiceAccountInvalid = err.NewErrorInvalidArgNotMatchFormat(
		err.PushPackageName,
		"fcm service account must be a service_account key file",
		"fcmServiceAccount",
	)
)

// ValidateFCMServiceAccount ensures the provided JSON is a Google
// service-account key file (type == "service_account") whose token endpoint is
// Google's. Any other credential type (notably "external_account",
// "external_account_authorized_user", "authorized_user" and
// "impersonated_service_account") is rejected to prevent credential-source
// abuse (arbitrary file read / SSRF) when the credential is later exchanged for
// an access token.
func ValidateFCMServiceAccount(fcmServiceAccount []byte) error {
	if len(fcmServiceAccount) == 0 {
		return ErrFCMServiceAccountRequired
	}
	var sa struct {
		Type      string          `json:"type"`
		TokenURI  string          `json:"token_uri"`
		Web       json.RawMessage `json:"web"`
		Installed json.RawMessage `json:"installed"`
	}
	if e := json.Unmarshal(fcmServiceAccount, &sa); e != nil {
		return ErrFCMServiceAccountInvalid
	}
	if sa.Type != "service_account" {
		return ErrFCMServiceAccountInvalid
	}
	// Reject OAuth client ("web"/"installed") blocks. google.CredentialsFromJSON
	// tries ConfigFromJSON first and prioritises such a block over the declared
	// type, returning a token source with no auth handler whose Token() call
	// panics in the subscriber. A service-account key file never contains them.
	if len(sa.Web) > 0 || len(sa.Installed) > 0 {
		return ErrFCMServiceAccountInvalid
	}
	if _, ok := allowedFCMTokenURIs[sa.TokenURI]; !ok {
		return ErrFCMServiceAccountInvalid
	}
	return nil
}
