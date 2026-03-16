// Copyright 2026 The Bucketeer Authors.
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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
)

func TestGetSuggestionsForPage(t *testing.T) {
	t.Parallel()

	t.Run("returns suggestions for feature flags page", func(t *testing.T) {
		t.Parallel()
		suggestions := getSuggestionsForPage(&aichatproto.PageContext{
			PageType: aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS,
		})
		require.Len(t, suggestions, 2)
		assert.Equal(t, "Organize flags with tags", suggestions[0].Title)
		assert.Equal(t, aichatproto.SuggestionType_SUGGESTION_TYPE_BEST_PRACTICE, suggestions[0].Type)
	})

	t.Run("returns suggestions for targeting page", func(t *testing.T) {
		t.Parallel()
		suggestions := getSuggestionsForPage(&aichatproto.PageContext{
			PageType: aichatproto.PageContext_PAGE_TYPE_TARGETING,
		})
		require.Len(t, suggestions, 2)
		assert.Equal(t, aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY, suggestions[0].Type)
	})

	t.Run("returns suggestions for experiments page", func(t *testing.T) {
		t.Parallel()
		suggestions := getSuggestionsForPage(&aichatproto.PageContext{
			PageType: aichatproto.PageContext_PAGE_TYPE_EXPERIMENTS,
		})
		require.Len(t, suggestions, 1)
	})

	t.Run("returns suggestions for autoops page", func(t *testing.T) {
		t.Parallel()
		suggestions := getSuggestionsForPage(&aichatproto.PageContext{
			PageType: aichatproto.PageContext_PAGE_TYPE_AUTOOPS,
		})
		require.Len(t, suggestions, 2)
	})

	t.Run("returns default suggestions for unspecified page", func(t *testing.T) {
		t.Parallel()
		suggestions := getSuggestionsForPage(&aichatproto.PageContext{
			PageType: aichatproto.PageContext_PAGE_TYPE_UNSPECIFIED,
		})
		require.Len(t, suggestions, 1)
	})

	t.Run("returns nil for nil page context", func(t *testing.T) {
		t.Parallel()
		suggestions := getSuggestionsForPage(nil)
		assert.Nil(t, suggestions)
	})

	t.Run("all suggestions have required fields", func(t *testing.T) {
		t.Parallel()
		pageTypes := []aichatproto.PageContext_PageType{
			aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS,
			aichatproto.PageContext_PAGE_TYPE_TARGETING,
			aichatproto.PageContext_PAGE_TYPE_EXPERIMENTS,
			aichatproto.PageContext_PAGE_TYPE_AUTOOPS,
			aichatproto.PageContext_PAGE_TYPE_UNSPECIFIED,
		}
		for _, pt := range pageTypes {
			suggestions := getSuggestionsForPage(&aichatproto.PageContext{PageType: pt})
			for _, s := range suggestions {
				assert.NotEmpty(t, s.Id)
				assert.NotEmpty(t, s.Title)
				assert.NotEmpty(t, s.Description)
				assert.NotEmpty(t, s.DocUrl)
				assert.NotEqual(t, aichatproto.SuggestionType_SUGGESTION_TYPE_UNSPECIFIED, s.Type)
			}
		}
	})
}
