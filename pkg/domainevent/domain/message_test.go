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

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	proto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestLocalizedMessage(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		inputEventType proto.Event_Type
		expected       *proto.LocalizedMessage
	}{
		"unknown match": {
			inputEventType: proto.Event_UNKNOWN,
			expected: &proto.LocalizedMessage{
				Locale:  locale.JaJP,
				Message: "不明な操作を実行しました",
			},
		},
		"unmatch": {
			inputEventType: proto.Event_Type(-1),
			expected: &proto.LocalizedMessage{
				Locale:  locale.JaJP,
				Message: "不明な操作を実行しました",
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := LocalizedMessage(p.inputEventType, locale.JaJP)
			assert.Equal(t, p.expected, actual)
		})
	}
}

// TestImplementedLocalizedMessage checks if every domain event type has a message.
func TestImplementedLocalizedMessage(t *testing.T) {
	t.Parallel()
	unknown := &proto.LocalizedMessage{
		Locale:  locale.JaJP,
		Message: "不明な操作を実行しました",
	}
	for k, v := range proto.Event_Type_name {
		if v == "UNKNOWN" {
			continue
		}
		t.Run(v, func(t *testing.T) {
			actual := LocalizedMessage(proto.Event_Type(k), locale.JaJP)
			assert.NotEqual(t, unknown, actual)
		})
	}
}
