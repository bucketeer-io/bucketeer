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
	"fmt"
	"strings"
	"unicode/utf8"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// buildFeatureContext creates a privacy-safe text representation of a feature flag
// for inclusion in the AI chat system prompt.
// Per RFC 0045, we include: flag name, description, variation names/descriptions,
// tag names, rule structure (operator, strategy type).
// We exclude: attribute values, user IDs, variation values, clause values.
// maxFeatureContextLength is the maximum byte length for the feature context text
// to avoid excessive LLM token consumption.
const maxFeatureContextLength = 2000

func buildFeatureContext(f *featureproto.Feature) string {
	if f == nil {
		return ""
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "Name: %s\n", f.Name)
	if f.Description != "" {
		fmt.Fprintf(&sb, "Description: %s\n", f.Description)
	}
	fmt.Fprintf(&sb, "Enabled: %t\n", f.Enabled)
	fmt.Fprintf(&sb, "Variation Type: %s\n", f.VariationType.String())

	// Variations — names and descriptions only; values are excluded for privacy.
	if len(f.Variations) > 0 {
		sb.WriteString("Variations:\n")
		for _, v := range f.Variations {
			name := v.Name
			if name == "" {
				name = v.Id
			}
			sb.WriteString("  - " + name)
			if v.Description != "" {
				fmt.Fprintf(&sb, " (%s)", v.Description)
			}
			sb.WriteString("\n")
		}
	}

	// Tags
	if len(f.Tags) > 0 {
		fmt.Fprintf(&sb, "Tags: %s\n", strings.Join(f.Tags, ", "))
	}

	// Rules — structure only; clause values and attribute names are excluded for privacy.
	if len(f.Rules) > 0 {
		fmt.Fprintf(&sb, "Targeting Rules: %d rule(s)\n", len(f.Rules))
		for i, rule := range f.Rules {
			fmt.Fprintf(&sb, "  Rule %d:\n", i+1)
			if rule.Strategy != nil {
				fmt.Fprintf(&sb, "    Strategy: %s\n", rule.Strategy.Type.String())
			}
			fmt.Fprintf(&sb, "    Conditions: %d\n", len(rule.Clauses))
			for _, clause := range rule.Clauses {
				fmt.Fprintf(&sb, "      - Operator: %s\n", clause.Operator.String())
			}
		}
	}

	// Default strategy
	if f.DefaultStrategy != nil {
		fmt.Fprintf(&sb, "Default Strategy: %s\n", f.DefaultStrategy.Type.String())
	}

	// Prerequisites
	if len(f.Prerequisites) > 0 {
		sb.WriteString("Prerequisites:\n")
		for _, p := range f.Prerequisites {
			fmt.Fprintf(&sb, "  - Depends on flag: %s\n", p.FeatureId)
		}
	}

	result := sb.String()
	if utf8.RuneCountInString(result) > maxFeatureContextLength {
		runes := []rune(result)
		result = string(runes[:maxFeatureContextLength]) + "\n... (truncated)\n"
	}
	return result
}
