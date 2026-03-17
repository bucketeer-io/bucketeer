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

package rag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripMDX(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		input       string
		contains    []string
		notContains []string
		exact       string
	}{
		{
			desc:        "removes import statements",
			input:       "import Tabs from '@theme/Tabs';\nimport TabItem from '@theme/TabItem';\n\n# Hello\nContent here.",
			contains:    []string{"# Hello", "Content here."},
			notContains: []string{"import"},
		},
		{
			desc:        "removes JSX self-closing tags",
			input:       "Some text\n<CustomComponent prop=\"value\" />\nMore text",
			contains:    []string{"Some text", "More text"},
			notContains: []string{"CustomComponent"},
		},
		{
			desc:        "removes JSX block tags",
			input:       "Before\n<Tabs>\n<TabItem value=\"go\" label=\"Go\">\n\n```go\nfmt.Println(\"hello\")\n```\n\n</TabItem>\n</Tabs>\nAfter",
			contains:    []string{"fmt.Println", "After"},
			notContains: []string{"<Tabs>", "<TabItem", "</TabItem>", "</Tabs>"},
		},
		{
			desc:        "removes export statements",
			input:       "export const meta = { title: 'Test' };\n\n# Title\nContent",
			contains:    []string{"# Title"},
			notContains: []string{"export"},
		},
		{
			desc:        "removes HTML comments",
			input:       "Text before\n<!-- This is a comment -->\nText after",
			contains:    []string{"Text before", "Text after"},
			notContains: []string{"comment"},
		},
		{
			desc:     "preserves code blocks",
			input:    "# Example\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n```\n\nEnd.",
			contains: []string{"package main", "fmt.Println"},
		},
		{
			desc:        "collapses multiple blank lines",
			input:       "Line 1\n\n\n\n\nLine 2",
			notContains: []string{"\n\n\n"},
		},
		{
			desc:  "handles empty input",
			input: "",
			exact: "",
		},
		{
			desc:        "preserves import inside code block",
			input:       "import Tabs from '@theme/Tabs';\n\n# Example\n\n```javascript\nimport React from 'react';\nconsole.log('hello');\n```\n\nEnd.",
			contains:    []string{"import React from 'react';", "console.log"},
			notContains: []string{"@theme/Tabs"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := StripMDX(p.input)
			if p.exact != "" || (p.input == "" && len(p.contains) == 0) {
				assert.Equal(t, p.exact, result)
				return
			}
			for _, s := range p.contains {
				assert.Contains(t, result, s)
			}
			for _, s := range p.notContains {
				assert.NotContains(t, result, s)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "extracts from frontmatter",
			input:    "---\ntitle: Feature Flags\nslug: /feature-flags\n---\n\n# Feature Flags\nContent",
			expected: "Feature Flags",
		},
		{
			desc:     "extracts from first h1 when no frontmatter",
			input:    "# Getting Started\n\nWelcome to Bucketeer.",
			expected: "Getting Started",
		},
		{
			desc:     "returns empty when no title found",
			input:    "Just some text without headers.",
			expected: "",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			title := ExtractTitle(p.input)
			assert.Equal(t, p.expected, title)
		})
	}
}
