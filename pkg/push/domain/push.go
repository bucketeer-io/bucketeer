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
	"errors"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/push"
)

var (
	ErrTagDuplicated    = errors.New("push: tag is duplicated")
	ErrTagAlreadyExists = errors.New("push: tag already exists")
	ErrTagNotFound      = errors.New("push: tag not found")
)

type Push struct {
	*proto.Push
}

func NewPush(name, fcmAPIKey string, tags []string) (*Push, error) {
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
		Name:      name,
		Id:        id.String(),
		FcmApiKey: fcmAPIKey,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}}
	return p, nil
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
