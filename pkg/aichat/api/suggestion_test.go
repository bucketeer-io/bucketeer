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

	patterns := []struct {
		desc            string
		pageContext     *aichatproto.PageContext
		expectedLen     int
		expectedNil     bool
		checkFirstTitle string
		checkFirstType  aichatproto.SuggestionType
	}{
		{
			desc: "returns suggestions for feature flags page",
			pageContext: &aichatproto.PageContext{
				PageType: aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS,
			},
			expectedLen:     2,
			checkFirstTitle: "Organize flags with tags",
			checkFirstType:  aichatproto.SuggestionType_SUGGESTION_TYPE_BEST_PRACTICE,
		},
		{
			desc: "returns suggestions for targeting page",
			pageContext: &aichatproto.PageContext{
				PageType: aichatproto.PageContext_PAGE_TYPE_TARGETING,
			},
			expectedLen:    2,
			checkFirstType: aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY,
		},
		{
			desc: "returns suggestions for experiments page",
			pageContext: &aichatproto.PageContext{
				PageType: aichatproto.PageContext_PAGE_TYPE_EXPERIMENTS,
			},
			expectedLen: 1,
		},
		{
			desc: "returns suggestions for autoops page",
			pageContext: &aichatproto.PageContext{
				PageType: aichatproto.PageContext_PAGE_TYPE_AUTOOPS,
			},
			expectedLen: 2,
		},
		{
			desc: "returns default suggestions for unspecified page",
			pageContext: &aichatproto.PageContext{
				PageType: aichatproto.PageContext_PAGE_TYPE_UNSPECIFIED,
			},
			expectedLen: 1,
		},
		{
			desc:        "returns nil for nil page context",
			pageContext: nil,
			expectedNil: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			suggestions := getSuggestionsForPage(p.pageContext)
			if p.expectedNil {
				assert.Nil(t, suggestions)
				return
			}
			require.Len(t, suggestions, p.expectedLen)
			if p.checkFirstTitle != "" {
				assert.Equal(t, p.checkFirstTitle, suggestions[0].Title)
			}
			if p.checkFirstType != 0 {
				assert.Equal(t, p.checkFirstType, suggestions[0].Type)
			}
		})
	}

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
