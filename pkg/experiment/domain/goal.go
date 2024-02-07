// Copyright 2024 The Bucketeer Authors.
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

	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

type Goal struct {
	*proto.Goal
}

func NewGoal(id, name, description string) (*Goal, error) {
	now := time.Now().Unix()
	return &Goal{&proto.Goal{
		Id:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}}, nil
}

func (g *Goal) Rename(name string) error {
	g.Goal.Name = name
	g.Goal.UpdatedAt = time.Now().Unix()
	return nil
}

func (g *Goal) ChangeDescription(description string) error {
	g.Goal.Description = description
	g.Goal.UpdatedAt = time.Now().Unix()
	return nil
}

func (g *Goal) SetArchived() error {
	g.Goal.Archived = true
	g.Goal.UpdatedAt = time.Now().Unix()
	return nil
}

func (g *Goal) SetDeleted() error {
	g.Goal.Deleted = true
	g.Goal.UpdatedAt = time.Now().Unix()
	return nil
}
