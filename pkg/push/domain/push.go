// Copyright 2025 The Bucketeer Authors.
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
	"slices"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/wrapperspb"

	err "github.com/bucketeer-io/bucketeer/pkg/error"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/push"
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
