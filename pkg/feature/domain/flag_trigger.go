// Copyright 2023 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type FlagTrigger struct {
	*proto.FlagTrigger
}

func NewFlagTrigger(
	environmentNamespace string,
	cmd *proto.CreateFlagTriggerCommand,
) (*FlagTrigger, error) {
	now := time.Now().Unix()
	triggerID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &FlagTrigger{&proto.FlagTrigger{
		Id:                   triggerID.String(),
		FeatureId:            cmd.FeatureId,
		EnvironmentNamespace: environmentNamespace,
		Type:                 cmd.Type,
		Action:               cmd.Action,
		Description:          cmd.Description,
		Disabled:             false,
		CreatedAt:            now,
		UpdatedAt:            now,
	}}, nil
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
	secret := NewFlagTriggerSecret(
		ft.GetId(),
		ft.GetFeatureId(),
		ft.GetEnvironmentNamespace(),
		newTriggerUuid.String(),
		int(ft.GetAction()),
	)
	encoded, err := secret.Marshal()
	if err != nil {
		return err
	}
	h := sha256.New()
	h.Write(encoded)
	hashed := h.Sum(nil)
	ft.Token = base64.RawURLEncoding.EncodeToString(hashed)
	return nil
}
