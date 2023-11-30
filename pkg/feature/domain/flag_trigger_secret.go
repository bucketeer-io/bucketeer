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

import "encoding/json"

type FlagTriggerSecret struct {
	ID                   string `json:"id"`
	FeatureID            string `json:"feature_id"`
	Action               int    `json:"action"`
	EnvironmentNamespace string `json:"environment_namespace"`
	UUID                 string `json:"uuid"`
}

func NewFlagTriggerSecret(id, featureID, namespace, uuid string, action int) *FlagTriggerSecret {
	return &FlagTriggerSecret{
		ID:                   id,
		FeatureID:            featureID,
		EnvironmentNamespace: namespace,
		UUID:                 uuid,
		Action:               action,
	}
}

func (f *FlagTriggerSecret) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FlagTriggerSecret) GetID() string {
	return f.ID
}

func (f *FlagTriggerSecret) GetFeatureID() string {
	return f.FeatureID
}

func (f *FlagTriggerSecret) GetAction() int {
	return f.Action
}

func (f *FlagTriggerSecret) GetUUID() string {
	return f.UUID
}

func (f *FlagTriggerSecret) GetEnvironmentNamespace() string {
	return f.EnvironmentNamespace
}

func UnmarshalFlagTriggerSecret(data []byte) (*FlagTriggerSecret, error) {
	f := FlagTriggerSecret{}
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}
