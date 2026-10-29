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
//

package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/jinzhu/copier"

	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type FlagTrigger struct {
	*proto.FlagTrigger
}

func NewFlagTrigger(
	environmentId string,
	featureId string,
	flagType proto.FlagTrigger_Type,
	action proto.FlagTrigger_Action,
	description string,
) (*FlagTrigger, error) {
	now := time.Now().Unix()
	triggerID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &FlagTrigger{&proto.FlagTrigger{
		Id:            triggerID.String(),
		FeatureId:     featureId,
		EnvironmentId: environmentId,
		Type:          flagType,
		Action:        action,
		Description:   description,
		Disabled:      false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}}, nil
}

func (ft *FlagTrigger) UpdateFlagTrigger(
	description *wrapperspb.StringValue,
	reset bool,
	disabled *wrapperspb.BoolValue,
) (*FlagTrigger, error) {
	updated := &FlagTrigger{}
	if err := copier.Copy(updated, ft); err != nil {
		return nil, err
	}
	if description != nil {
		updated.Description = description.Value
	}
	if reset {
		err := updated.GenerateToken()
		if err != nil {
			return nil, err
		}
	}
	if disabled != nil {
		updated.Disabled = disabled.Value
	}
	updated.UpdatedAt = time.Now().Unix()
	return updated, nil
}

func (ft *FlagTrigger) ChangeDescription(description string) error {
	ft.Description = description
	ft.UpdatedAt = time.Now().Unix()
	return nil
}

func (ft *FlagTrigger) Disable() error {
	ft.Disabled = true
	ft.UpdatedAt = time.Now().Unix()
	return nil
}

func (ft *FlagTrigger) Enable() error {
	ft.Disabled = false
	ft.UpdatedAt = time.Now().Unix()
	return nil
}

func (ft *FlagTrigger) UpdateTriggerUsage() error {
	unix := time.Now().Unix()
	ft.LastTriggeredAt = unix
	ft.UpdatedAt = unix
	ft.TriggerCount = ft.TriggerCount + 1
	return nil
}

func (ft *FlagTrigger) GenerateToken() error {
	newTriggerUuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	h := sha256.New()
	h.Write([]byte(newTriggerUuid.String()))
	hashed := h.Sum(nil)
	ft.Token = base64.RawURLEncoding.EncodeToString(hashed)
	return nil
}

// WithoutToken returns a deep copy of the flag trigger with an empty token.
// The token authorizes the trigger webhook, so it is returned only once, to
// the editor who created or reset the trigger. Every other record of the
// trigger, such as the domain event and the audit log built from it, must use
// this copy so the token is never stored outside the flag trigger table.
func (ft *FlagTrigger) WithoutToken() *FlagTrigger {
	copied := &FlagTrigger{FlagTrigger: pb.Clone(ft.FlagTrigger).(*proto.FlagTrigger)}
	copied.Token = ""
	return copied
}
