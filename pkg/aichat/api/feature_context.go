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
// All user-controlled fields are sanitized to mitigate prompt injection.
const maxFeatureContextLength = 2000
const maxFieldLength = 200

func buildFeatureContext(f *featureproto.Feature) string {
	if f == nil {
		return ""
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "Name: %q\n", sanitizePromptField(f.Name))
	if f.Description != "" {
		fmt.Fprintf(&sb, "Description: %q\n", sanitizePromptField(f.Description))
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
			fmt.Fprintf(&sb, "  - %q", sanitizePromptField(name))
			if v.Description != "" {
				fmt.Fprintf(&sb, " (%q)", sanitizePromptField(v.Description))
			}
			sb.WriteString("\n")
		}
	}

	// Tags
	if len(f.Tags) > 0 {
		sanitized := make([]string, len(f.Tags))
		for i, t := range f.Tags {
			sanitized[i] = sanitizePromptField(t)
		}
		fmt.Fprintf(&sb, "Tags: %s\n", strings.Join(sanitized, ", "))
	}

	// Rules — structure only; clause values and attribute names are excluded.
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
			featureID := sanitizePromptField(p.FeatureId)
			fmt.Fprintf(&sb, "  - Depends on flag: %q\n", featureID)
		}
	}

	result := sb.String()
	if utf8.RuneCountInString(result) > maxFeatureContextLength {
		runes := []rune(result)
		result = string(runes[:maxFeatureContextLength]) + "\n... (truncated)\n"
	}
	return "<feature_data>\n" + result + "</feature_data>"
}

// sanitizePromptField sanitizes a user-controlled string before embedding
// it in the system prompt. It removes control characters and newlines,
// collapses whitespace, and truncates to maxFieldLength runes.
func sanitizePromptField(s string) string {
	// Remove control characters (including \n, \r, \t)
	cleaned := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, s)
	// Collapse multiple spaces
	parts := strings.Fields(cleaned)
	cleaned = strings.Join(parts, " ")
	// Truncate
	if utf8.RuneCountInString(cleaned) > maxFieldLength {
		cleaned = string([]rune(cleaned)[:maxFieldLength])
	}
	return cleaned
}
