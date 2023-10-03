// Copyright 2023 The Bucketeer Authors.
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

package locale

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestMustLocalizeWithTemplate(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		fields   []string
		lang     string
		expected string
	}{
		{
			name:     "succeed",
			id:       RequiredFieldTemplate,
			fields:   []string{"field-1"},
			lang:     "",
			expected: "field-1 is required",
		},
		{
			name:     "succeed",
			id:       RequiredFieldTemplate,
			fields:   []string{"field-1"},
			lang:     Ja,
			expected: "field-1は必須です",
		},
		{
			name:     "succeed",
			id:       RequiredFieldTemplate,
			fields:   []string{"field-1"},
			lang:     En,
			expected: "field-1 is required",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{c.lang},
			})
			loc := NewLocalizer(ctx)
			actual := loc.MustLocalizeWithTemplate(c.id, c.fields...)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestMustLocalize(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		lang     string
		expected string
	}{
		{
			name:     "succeed",
			id:       FeatureFlagID,
			lang:     "",
			expected: "Feature flag ID",
		},
		{
			name:     "succeed",
			id:       FeatureFlagID,
			lang:     Ja,
			expected: "フィーチャーフラグID",
		},
		{
			name:     "succeed",
			id:       FeatureFlagID,
			lang:     En,
			expected: "Feature flag ID",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{c.lang},
			})
			loc := NewLocalizer(ctx)
			actual := loc.MustLocalize(c.id)
			assert.Equal(t, c.expected, actual)
		})
	}
}
