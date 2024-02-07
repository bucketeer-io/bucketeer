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

package uuid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrIncorrectUUIDFormat = errors.New("uuid: format must be an uuid version 4")
	// Version 4
	uuidRegex = regexp.MustCompile(
		"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$",
	)
)

// https://tools.ietf.org/html/rfc4122
type UUID [16]byte

func NewUUID() (*UUID, error) {
	uuid := &UUID{}
	if _, err := rand.Read(uuid[:]); err != nil {
		return nil, err
	}
	uuid.setVariant()
	uuid.setVersion()
	return uuid, nil
}

func ValidateUUID(id string) error {
	if !uuidRegex.MatchString(id) {
		return ErrIncorrectUUIDFormat
	}
	return nil
}

func (uuid *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (uuid *UUID) setVariant() {
	uuid[8] = (uuid[8] & 0x3f) | 0x80
}

func (uuid *UUID) setVersion() {
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
}
