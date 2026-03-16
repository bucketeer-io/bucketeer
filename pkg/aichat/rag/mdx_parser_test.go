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

	t.Run("removes import statements", func(t *testing.T) {
		t.Parallel()
		input := "import Tabs from '@theme/Tabs';\nimport TabItem from '@theme/TabItem';\n\n# Hello\nContent here."
		result := StripMDX(input)
		assert.NotContains(t, result, "import")
		assert.Contains(t, result, "# Hello")
		assert.Contains(t, result, "Content here.")
	})

	t.Run("removes JSX self-closing tags", func(t *testing.T) {
		t.Parallel()
		input := "Some text\n<CustomComponent prop=\"value\" />\nMore text"
		result := StripMDX(input)
		assert.NotContains(t, result, "CustomComponent")
		assert.Contains(t, result, "Some text")
		assert.Contains(t, result, "More text")
	})

	t.Run("removes JSX block tags", func(t *testing.T) {
		t.Parallel()
		input := "Before\n<Tabs>\n<TabItem value=\"go\" label=\"Go\">\n\n```go\nfmt.Println(\"hello\")\n```\n\n</TabItem>\n</Tabs>\nAfter"
		result := StripMDX(input)
		assert.NotContains(t, result, "<Tabs>")
		assert.NotContains(t, result, "<TabItem")
		assert.NotContains(t, result, "</TabItem>")
		assert.NotContains(t, result, "</Tabs>")
		assert.Contains(t, result, "fmt.Println")
		assert.Contains(t, result, "After")
	})

	t.Run("removes export statements", func(t *testing.T) {
		t.Parallel()
		input := "export const meta = { title: 'Test' };\n\n# Title\nContent"
		result := StripMDX(input)
		assert.NotContains(t, result, "export")
		assert.Contains(t, result, "# Title")
	})

	t.Run("removes HTML comments", func(t *testing.T) {
		t.Parallel()
		input := "Text before\n<!-- This is a comment -->\nText after"
		result := StripMDX(input)
		assert.NotContains(t, result, "comment")
		assert.Contains(t, result, "Text before")
		assert.Contains(t, result, "Text after")
	})

	t.Run("preserves code blocks", func(t *testing.T) {
		t.Parallel()
		input := "# Example\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n```\n\nEnd."
		result := StripMDX(input)
		assert.Contains(t, result, "package main")
		assert.Contains(t, result, "fmt.Println")
	})

	t.Run("collapses multiple blank lines", func(t *testing.T) {
		t.Parallel()
		input := "Line 1\n\n\n\n\nLine 2"
		result := StripMDX(input)
		assert.NotContains(t, result, "\n\n\n")
	})

	t.Run("handles empty input", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "", StripMDX(""))
	})

	t.Run("preserves import inside code block", func(t *testing.T) {
		t.Parallel()
		input := "import Tabs from '@theme/Tabs';\n\n# Example\n\n```javascript\nimport React from 'react';\nconsole.log('hello');\n```\n\nEnd."
		result := StripMDX(input)
		// MDX import should be removed
		assert.NotContains(t, result, "@theme/Tabs")
		// Code block import should be preserved
		assert.Contains(t, result, "import React from 'react';")
		assert.Contains(t, result, "console.log")
	})
}

func TestExtractTitle(t *testing.T) {
	t.Parallel()

	t.Run("extracts from frontmatter", func(t *testing.T) {
		t.Parallel()
		input := "---\ntitle: Feature Flags\nslug: /feature-flags\n---\n\n# Feature Flags\nContent"
		title := ExtractTitle(input)
		assert.Equal(t, "Feature Flags", title)
	})

	t.Run("extracts from first h1 when no frontmatter", func(t *testing.T) {
		t.Parallel()
		input := "# Getting Started\n\nWelcome to Bucketeer."
		title := ExtractTitle(input)
		assert.Equal(t, "Getting Started", title)
	})

	t.Run("returns empty when no title found", func(t *testing.T) {
		t.Parallel()
		input := "Just some text without headers."
		title := ExtractTitle(input)
		assert.Equal(t, "", title)
	})
}
