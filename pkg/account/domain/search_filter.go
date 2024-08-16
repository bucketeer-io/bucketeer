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
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

type SearchFilter struct {
	*proto.SearchFilter
}

func NewSearchFilter(
	name string,
	query string,
	defaultFilter bool,
	filterTargetType proto.FilterTargetType,
	environmentID string,
) (*SearchFilter, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &SearchFilter{&proto.SearchFilter{
		Id:               id.String(),
		Name:             name,
		Query:            query,
		DefaultFilter:    defaultFilter,
		FilterTargetType: filterTargetType,
		EnvironmentId:    environmentID,
	}}, nil
}

func (s *SearchFilter) UpdateSearchFilter(
	name string,
	query string,
	defaultFilter bool,
) {
	s.Name = name
	s.Query = query
	s.DefaultFilter = defaultFilter
}

func (s *SearchFilter) SetDefaultFilter(isDefault bool) {
	s.DefaultFilter = isDefault
}
