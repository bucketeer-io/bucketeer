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

package email

import (
	"strings"
	"testing"
)

func TestTemplateRenderer_RenderWelcomeEmail(t *testing.T) {
	t.Parallel()

	config := Config{
		BaseURL: "https://example.com",
	}
	renderer := NewTemplateRenderer(config)

	subject, body := renderer.RenderWelcomeEmail("en", "test@example.com")

	if subject == "" {
		t.Error("expected non-empty subject")
	}

	if !strings.Contains(body, "test@example.com") {
		t.Error("expected body to contain user email")
	}

	if !strings.Contains(body, "https://example.com") {
		t.Error("expected body to contain base URL")
	}

	if !strings.Contains(body, "invited to join") {
		t.Error("expected body to contain welcome message")
	}
}

func TestTemplateRenderer_LanguageFallback(t *testing.T) {
	t.Parallel()

	config := Config{
		BaseURL: "https://example.com",
		Templates: TemplatesConfig{
			En: TemplatesByLanguage{
				Welcome: Template{
					Subject: "Welcome to Bucketeer",
					Body:    "Welcome {{userEmail}}!",
				},
			},
			Ja: TemplatesByLanguage{
				Welcome: Template{
					Subject: "Bucketeerへようこそ",
					Body:    "ようこそ {{userEmail}}!",
				},
			},
		},
	}
	renderer := NewTemplateRenderer(config)

	// Test English
	subjectEn, bodyEn := renderer.RenderWelcomeEmail("en", "test@example.com")
	if subjectEn != "Welcome to Bucketeer" {
		t.Errorf("expected English subject, got: %s", subjectEn)
	}
	if !strings.Contains(bodyEn, "Welcome test@example.com") {
		t.Errorf("expected English body to contain 'Welcome test@example.com', got: %s", bodyEn)
	}

	// Test Japanese
	subjectJa, bodyJa := renderer.RenderWelcomeEmail("ja", "test@example.com")
	if subjectJa != "Bucketeerへようこそ" {
		t.Errorf("expected Japanese subject, got: %s", subjectJa)
	}
	if !strings.Contains(bodyJa, "ようこそ test@example.com") {
		t.Errorf("expected Japanese body to contain 'ようこそ test@example.com', got: %s", bodyJa)
	}

	// Test unsupported language falls back to English
	subjectFallback, _ := renderer.RenderWelcomeEmail("fr", "test@example.com")
	if subjectFallback != "Welcome to Bucketeer" {
		t.Errorf("expected fallback to English subject, got: %s", subjectFallback)
	}
}
