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
	"time"

	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

// Webhook holds the settings for accepting webhooks from alert systems, etc.
type Webhook struct {
	*proto.Webhook
}

func NewWebhook(id, name, description string) *Webhook {
	now := time.Now().Unix()
	return &Webhook{&proto.Webhook{
		Id:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}}
}

func (w *Webhook) ChangeName(name string) error {
	w.Name = name
	w.UpdatedAt = time.Now().Unix()
	return nil
}

func (w *Webhook) ChangeDescription(description string) error {
	w.Description = description
	w.UpdatedAt = time.Now().Unix()
	return nil
}
