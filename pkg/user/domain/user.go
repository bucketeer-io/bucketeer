// Copyright 2022 The Bucketeer Authors.
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
	"strings"

	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

var (
	ErrNotSameID = errors.New("user: user id is not same")
	ErrNotLater  = errors.New("user: update user is not later")
)

type User struct {
	*userproto.User
}

func (u *User) ID() string {
	return u.Id
}

func (u *User) UpdateMe(newer *User) error {
	if u.Id != newer.Id {
		return ErrNotSameID
	}
	if u.LastSeen >= newer.LastSeen {
		return ErrNotLater
	}
	u.LastSeen = newer.LastSeen
	for key, newData := range newer.TaggedData {
		u.TaggedData[key] = u.trimData(newData)
	}
	return nil
}

func (u *User) trimData(data *userproto.User_Data) *userproto.User_Data {
	if data == nil {
		return nil
	}
	for key, val := range data.Value {
		data.Value[key] = strings.TrimSpace(val)
	}
	return data
}

func (u *User) Data(tag string) map[string]string {
	if u.TaggedData[tag] == nil {
		return nil
	}
	return u.TaggedData[tag].Value
}
