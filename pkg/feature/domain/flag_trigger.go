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
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type FlagTrigger struct {
	*proto.FlagTrigger
}

func NewFlagTrigger(
	namespace string,
	cmd *proto.CreateFlagTriggerCommand,
) (*FlagTrigger, error) {
	now := time.Now().Unix()
	triggerID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	triggerUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &FlagTrigger{&proto.FlagTrigger{
		Id:                   triggerID.String(),
		FeatureId:            cmd.FeatureId,
		EnvironmentNamespace: namespace,
		Type:                 cmd.Type,
		Action:               cmd.Action,
		Description:          cmd.Description,
		Uuid:                 triggerUUID.String(),
		Disabled:             false,
		Deleted:              false,
		CreatedAt:            now,
		UpdatedAt:            now,
	}}, nil
}

func (ft *FlagTrigger) GetId() string {
	return ft.Id
}

func (ft *FlagTrigger) GetDescription() string {
	return ft.Description
}
