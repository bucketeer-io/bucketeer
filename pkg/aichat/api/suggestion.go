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
	"context"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
)

// GetSuggestions returns proactive suggestions based on page context.
func (s *AIChatService) GetSuggestions(
	ctx context.Context,
	req *aichatproto.GetSuggestionsRequest,
) (*aichatproto.GetSuggestionsResponse, error) {
	// Input validation (before authorization to avoid unnecessary RPC)
	if req.EnvironmentId == "" {
		return nil, statusMissingEnvironmentID.Err()
	}

	// Authorization check (Viewer role minimum)
	_, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId,
	)
	if err != nil {
		return nil, err
	}

	suggestions := getSuggestionsForPage(req.PageContext)
	return &aichatproto.GetSuggestionsResponse{
		Suggestions: suggestions,
	}, nil
}

// getSuggestionsForPage returns rule-based suggestions for the given page context.
// No LLM calls are required.
func getSuggestionsForPage(
	pageCtx *aichatproto.PageContext,
) []*aichatproto.Suggestion {
	if pageCtx == nil {
		return nil
	}

	switch pageCtx.PageType {
	case aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS:
		return []*aichatproto.Suggestion{
			{
				Id:          "sug-ff-001",
				Title:       "Organize flags with tags",
				Description: "Tags help organize flags and optimize SDK performance by fetching only needed flags.",
				DocUrl:      "https://docs.bucketeer.io/best-practices/optimize-with-tags",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_BEST_PRACTICE,
			},
			{
				Id:          "sug-ff-002",
				Title:       "Archive unused flags",
				Description: "Flags that have been fully rolled out or unused for a long time should be archived to reduce flag debt.",
				DocUrl:      "https://docs.bucketeer.io/best-practices/feature-flag-lifecycle",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_OPTIMIZATION,
			},
		}
	case aichatproto.PageContext_PAGE_TYPE_TARGETING:
		return []*aichatproto.Suggestion{
			{
				Id:          "sug-tgt-001",
				Title:       "Simplify with Segments",
				Description: "If you use the same targeting conditions across multiple flags, create a Segment to reuse them.",
				DocUrl:      "https://docs.bucketeer.io/feature-flags/segments",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY,
			},
			{
				Id:          "sug-tgt-002",
				Title:       "Use Prerequisite Rules",
				Description: "Set up dependencies between flags to coordinate feature releases across microservices.",
				DocUrl:      "https://docs.bucketeer.io/feature-flags/creating-feature-flags/prerequisite",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY,
			},
		}
	case aichatproto.PageContext_PAGE_TYPE_EXPERIMENTS:
		return []*aichatproto.Suggestion{
			{
				Id:          "sug-exp-001",
				Title:       "Set up Goals first",
				Description: "Define what you want to measure before starting an experiment. Goals track user actions for statistical analysis.",
				DocUrl:      "https://docs.bucketeer.io/feature-flags/testing-with-flags/experiments",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_BEST_PRACTICE,
			},
		}
	case aichatproto.PageContext_PAGE_TYPE_AUTOOPS:
		return []*aichatproto.Suggestion{
			{
				Id:          "sug-ao-001",
				Title:       "Try Progressive Rollout",
				Description: "Gradually increase the percentage of users seeing a feature to minimize risk during deployment.",
				DocUrl:      "https://docs.bucketeer.io/feature-flags/auto-operation/progressive-rollout",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY,
			},
			{
				Id:          "sug-ao-002",
				Title:       "Automate with Flag Triggers",
				Description: "Use webhooks to control flags from CI/CD pipelines or external monitoring systems.",
				DocUrl:      "https://docs.bucketeer.io/feature-flags/flag-triggers",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY,
			},
		}
	default:
		return []*aichatproto.Suggestion{
			{
				Id:          "sug-gen-001",
				Title:       "Explore Bucketeer's features",
				Description: "Bucketeer offers feature flags, A/B testing, segments, progressive rollout, and more.",
				DocUrl:      "https://docs.bucketeer.io",
				Type:        aichatproto.SuggestionType_SUGGESTION_TYPE_FEATURE_DISCOVERY,
			},
		}
	}
}
