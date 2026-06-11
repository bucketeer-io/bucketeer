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
	"fmt"
	"regexp"
	"strings"
)

var (
	// Matches fenced code blocks: ```...```
	reCodeBlock = regexp.MustCompile("(?s)```.*?```")
	// Matches import statements: import X from 'Y';
	reImport = regexp.MustCompile(`(?m)^import\s+.*;\s*$`)
	// Matches export statements: export const/default ...
	reExport = regexp.MustCompile(`(?m)^export\s+.*$`)
	// Matches self-closing JSX tags: <Component prop="value" />
	reJSXSelfClosing = regexp.MustCompile(`<[A-Z][A-Za-z]*(?:\s[^>]*)?\s*/>`)
	// Matches JSX opening tags: <Component ...>
	reJSXOpen = regexp.MustCompile(`<[A-Z][A-Za-z]*(?:\s[^>]*)?>`)
	// Matches JSX closing tags: </Component>
	reJSXClose = regexp.MustCompile(`</[A-Z][A-Za-z]*>`)
	// Matches HTML comments: <!-- ... -->
	reHTMLComment = regexp.MustCompile(`<!--[\s\S]*?-->`)
	// Matches 3+ consecutive newlines
	reMultiNewline = regexp.MustCompile(`\n{3,}`)
	// Matches frontmatter block
	reFrontmatter = regexp.MustCompile(`(?s)^---\n.*?\n---\n?`)
	// Extracts title from frontmatter
	reFrontmatterTitle = regexp.MustCompile(`(?m)^title:\s*(.+)$`)
	// Matches first h1 header
	reH1 = regexp.MustCompile(`(?m)^#\s+(.+)$`)
)

// StripMDX removes MDX/JSX syntax from markdown content, leaving clean text.
// Code blocks are preserved as they contain useful SDK examples.
func StripMDX(input string) string {
	if input == "" {
		return ""
	}

	// Extract code blocks and replace with placeholders to protect them
	var codeBlocks []string
	result := reCodeBlock.ReplaceAllStringFunc(input, func(match string) string {
		idx := len(codeBlocks)
		codeBlocks = append(codeBlocks, match)
		return placeholderFor(idx)
	})

	// Remove frontmatter
	result = reFrontmatter.ReplaceAllString(result, "")
	// Remove import statements (outside code blocks)
	result = reImport.ReplaceAllString(result, "")
	// Remove export statements (outside code blocks)
	result = reExport.ReplaceAllString(result, "")
	// Remove HTML comments
	result = reHTMLComment.ReplaceAllString(result, "")
	// Remove self-closing JSX tags
	result = reJSXSelfClosing.ReplaceAllString(result, "")
	// Remove JSX opening tags
	result = reJSXOpen.ReplaceAllString(result, "")
	// Remove JSX closing tags
	result = reJSXClose.ReplaceAllString(result, "")

	// Restore code blocks in a single pass
	if len(codeBlocks) > 0 {
		oldnew := make([]string, 0, len(codeBlocks)*2)
		for i, block := range codeBlocks {
			oldnew = append(oldnew, placeholderFor(i), block)
		}
		result = strings.NewReplacer(oldnew...).Replace(result)
	}

	// Collapse multiple blank lines
	result = reMultiNewline.ReplaceAllString(result, "\n\n")

	return strings.TrimSpace(result)
}

func placeholderFor(idx int) string {
	return fmt.Sprintf("\x00CODEBLOCK_%d\x00", idx)
}

// ExtractTitle extracts the document title from frontmatter or first h1 header.
func ExtractTitle(input string) string {
	// Try frontmatter title first
	if matches := reFrontmatterTitle.FindStringSubmatch(input); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	// Fall back to first h1
	if matches := reH1.FindStringSubmatch(input); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
