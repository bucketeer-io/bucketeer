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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	proto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestLocalizedMessage(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	patterns := []struct {
		desc           string
		inputEventType proto.Event_Type
		expected       *proto.LocalizedMessage
	}{
		{
			desc:           "unknown match",
			inputEventType: proto.Event_UNKNOWN,
			expected: &proto.LocalizedMessage{
				Locale:  locale.Ja,
				Message: "不明な操作が実行されました",
			},
		},
		{
			desc:           "unmatch",
			inputEventType: proto.Event_Type(-1),
			expected: &proto.LocalizedMessage{
				Locale:  locale.Ja,
				Message: "不明な操作が実行されました",
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := LocalizedMessage(p.inputEventType, localizer)
			assert.Equal(t, p.expected, actual)
		})
	}
}

// TestImplementedLocalizedMessage checks if every domain event type has a message.
func TestImplementedLocalizedMessage(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	unknown := &proto.LocalizedMessage{
		Locale:  locale.Ja,
		Message: "不明な操作が実行されました",
	}
	for k, v := range proto.Event_Type_name {
		if v == "UNKNOWN" {
			continue
		}
		t.Run(v, func(t *testing.T) {
			actual := LocalizedMessage(proto.Event_Type(k), localizer)
			assert.NotEqual(t, unknown, actual)
		})
	}
}
