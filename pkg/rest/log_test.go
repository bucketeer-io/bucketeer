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

package rest

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBody(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		body        io.Reader
		expected    interface{}
		expectedErr bool
	}{
		{
			desc:        "err: not json",
			body:        strings.NewReader(`{tag: "ios", user: {id: "pingdom", data: {foo: "bar"}}}`),
			expected:    nil,
			expectedErr: true,
		},
		{
			desc:        "success: nil",
			body:        bytes.NewReader(nil),
			expected:    nil,
			expectedErr: false,
		},
		{
			desc: "success: json",
			body: strings.NewReader(`{"tag":"ios","user":{"id":"pingdom","data":{"foo":"bar"}}}`),
			expected: map[string]interface{}{
				"tag": "ios",
				"user": map[string]interface{}{
					"id": "pingdom",
					"data": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			decoded, err := decodeBody(p.body)
			assert.Equal(t, p.expected, decoded)
			assert.Equal(t, p.expectedErr, err != nil)
		})
	}
}
