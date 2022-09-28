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
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestUpdateMe(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		origin      *User
		input       *User
		expected    *User
		expectedErr error
	}{
		"update without data": {
			origin: &User{
				User: &userproto.User{
					Id:       "hoge",
					LastSeen: 0,
				},
			},
			input: &User{
				User: &userproto.User{
					Id:       "hoge",
					LastSeen: 1,
				},
			},
			expected: &User{
				User: &userproto.User{
					Id:       "hoge",
					Data:     map[string]string{},
					LastSeen: 1,
				},
			},
		},
		"update overriding data": {
			origin: &User{
				User: &userproto.User{
					Id:         "id",
					TaggedData: map[string]*userproto.User_Data{"tag": {Value: map[string]string{"key-0": "val-0", "key-1": "val-1"}}},
					LastSeen:   0,
				},
			},
			input: &User{
				User: &userproto.User{
					Id:         "id",
					TaggedData: map[string]*userproto.User_Data{"tag": {Value: map[string]string{"key-0": " val-0 ", "key-2": "  val-2  "}}},
					LastSeen:   1,
				},
			},
			expected: &User{
				User: &userproto.User{
					Id:         "id",
					TaggedData: map[string]*userproto.User_Data{"tag": {Value: map[string]string{"key-0": "val-0", "key-2": "val-2"}}},
					LastSeen:   1,
				},
			},
		},
		"update appending data": {
			origin: &User{
				User: &userproto.User{
					Id: "id",
					TaggedData: map[string]*userproto.User_Data{
						"tag-0": {Value: map[string]string{"key-0": "val-0", "key-1": "val-1"}},
						"tag-1": {Value: map[string]string{"key-1": "val-1", "key-2": "val-2"}},
					},
					LastSeen: 0,
				},
			},
			input: &User{
				User: &userproto.User{
					Id: "id",
					TaggedData: map[string]*userproto.User_Data{
						"tag-2": {Value: map[string]string{"key-2": " val-2 ", "key-3": "  val-3  "}},
						"tag-3": {Value: map[string]string{"key-3": " val-3 ", "key-4": "  val-4  "}},
					},
					LastSeen: 1,
				},
			},
			expected: &User{
				User: &userproto.User{
					Id: "id",
					TaggedData: map[string]*userproto.User_Data{
						"tag-0": {Value: map[string]string{"key-0": "val-0", "key-1": "val-1"}},
						"tag-1": {Value: map[string]string{"key-1": "val-1", "key-2": "val-2"}},
						"tag-2": {Value: map[string]string{"key-2": "val-2", "key-3": "val-3"}},
						"tag-3": {Value: map[string]string{"key-3": "val-3", "key-4": "val-4"}},
					},
					LastSeen: 1,
				},
			},
		},
		"err: id not same": {
			origin: &User{
				User: &userproto.User{
					Id:       "foo",
					LastSeen: 0,
				},
			},
			input: &User{
				User: &userproto.User{
					Id:       "fee",
					LastSeen: 1,
				},
			},
			expected: &User{
				User: &userproto.User{
					Id:       "foo",
					LastSeen: 0,
				},
			},
			expectedErr: ErrNotSameID,
		},
		"err: id not later": {
			origin: &User{
				User: &userproto.User{
					Id:       "foo",
					LastSeen: 1,
				},
			},
			input: &User{
				User: &userproto.User{
					Id:       "foo",
					LastSeen: 0,
				},
			},
			expected: &User{
				User: &userproto.User{
					Id:       "foo",
					LastSeen: 1,
				},
			},
			expectedErr: ErrNotLater,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			err := p.origin.UpdateMe(p.input)
			assert.True(t, proto.Equal(p.expected, p.origin))
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestData(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		origin   *User
		input    string
		expected map[string]string
	}{
		"no data": {
			origin: &User{
				User: &userproto.User{
					TaggedData: map[string]*userproto.User_Data{
						"t0": {Value: map[string]string{"t0-k0": "t0-v0"}},
					},
				},
			},
			input:    "t1",
			expected: nil,
		},
		"hit": {
			origin: &User{
				User: &userproto.User{
					TaggedData: map[string]*userproto.User_Data{
						"t0": {Value: map[string]string{"t0-k0": "t0-v0"}},
						"t1": {Value: map[string]string{"t1-k0": "t1-v0"}},
					},
				},
			},
			input: "t1",
			expected: map[string]string{
				"t1-k0": "t1-v0",
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := p.origin.Data(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}
