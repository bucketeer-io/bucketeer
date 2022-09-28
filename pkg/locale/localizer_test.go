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

package locale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustLocalizeWithTemplate(t *testing.T) {
	lJA := NewLocalizer(NewLocale(JaJP))
	lEN := NewLocalizer(NewLocale(EnUS))
	cases := []struct {
		name     string
		id       string
		fields   []string
		l        Localizer
		expected string
	}{
		{
			name:     "succeed",
			id:       RequiredFieldTemplate,
			fields:   []string{"field-1"},
			l:        lJA,
			expected: "field-1は必須です",
		},
		{
			name:     "succeed",
			id:       RequiredFieldTemplate,
			fields:   []string{"field-1"},
			l:        lEN,
			expected: "field-1 is required",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.l.MustLocalizeWithTemplate(c.id, c.fields...)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestMustLocalize(t *testing.T) {
	lJA := NewLocalizer(NewLocale(JaJP))
	lEN := NewLocalizer(NewLocale(EnUS))
	cases := []struct {
		name     string
		id       string
		l        Localizer
		expected string
	}{
		{
			name:     "succeed",
			id:       FeatureFlagID,
			l:        lJA,
			expected: "フィーチャーフラグID",
		},
		{
			name:     "succeed",
			id:       FeatureFlagID,
			l:        lEN,
			expected: "feature flag ID",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.l.MustLocalize(c.id)
			assert.Equal(t, c.expected, actual)
		})
	}
}
