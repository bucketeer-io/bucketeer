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

package domain

import (
	"encoding/json"
)

type webhookSecret struct {
	WebhookID            string `json:"webhook_id"`
	EnvironmentNamespace string `json:"environment_namespace"`
}

type WebhookSecret interface {
	Marshal() ([]byte, error)
	GetWebhookID() string
	GetEnvironmentNamespace() string
}

func NewWebhookSecret(webhookID, environmentNamespace string) WebhookSecret {
	return &webhookSecret{
		EnvironmentNamespace: environmentNamespace,
		WebhookID:            webhookID,
	}
}

func UnmarshalWebhookSecret(data []byte) (WebhookSecret, error) {
	ws := webhookSecret{}
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

func (ws *webhookSecret) Marshal() ([]byte, error) {
	return json.Marshal(ws)
}

func (ws *webhookSecret) GetWebhookID() string {
	return ws.WebhookID
}

func (ws *webhookSecret) GetEnvironmentNamespace() string {
	return ws.EnvironmentNamespace
}
