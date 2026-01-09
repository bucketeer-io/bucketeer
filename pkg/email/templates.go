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
)

// Default email templates
var (
	defaultWelcomeSubject = "You've been invited to join Bucketeer"
	//nolint:lll
	defaultWelcomeBody = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="x-apple-disable-message-reformatting">
    <title>You've been invited to join Bucketeer</title>
  </head>
  <body style="margin:0;padding:0;background:#f5f5f5;color:#333;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;line-height:1.6;">
    <div style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0;">
      Sign in with the authentication method configured for your account.
    </div>
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
      <tr>
        <td align="center" style="padding:24px 12px;">
          <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="width:100%;max-width:600px;background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
            <tr>
              <td style="padding:32px;background:#ffffff;">
                <img src="{{webConsoleEndpoint}}/img/bucketeer-logo-primary.png" alt="Bucketeer" width="205" height="48" style="display:block;margin-bottom:24px;" />
                <p style="margin:0 0 16px 0;font-size:15px;color:#4b5563;">Hello, {{userEmail}}.</p>
                <p style="margin:0 0 16px 0;font-size:15px;color:#4b5563;">
                  You've been invited to join <strong>Bucketeer</strong>, a feature flag and A/B testing platform.
                </p>
                <p style="margin:0 0 20px 0;font-size:15px;color:#4b5563;">
                  Use the link below to access the Bucketeer console and sign in with the authentication method configured for your account.
                </p>
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:0 0 20px 0;">
                  <tr><td align="center">
                    <a href="{{webConsoleEndpoint}}" style="background:#5B21B6;color:#ffffff;text-decoration:none;padding:14px 32px;border-radius:6px;display:inline-block;font-weight:600;font-size:15px;">Open Bucketeer Console</a>
                  </td></tr>
                </table>
                <p style="margin:0 0 8px 0;font-size:13px;color:#4b5563;">Or copy and paste this link into your browser:</p>
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:0 0 24px 0;">
                  <tr><td style="padding:12px 16px;background:#f9fafb;border:1px solid #e5e7eb;border-radius:6px;word-break:break-all;font-size:13px;color:#4b5563;">
                    {{webConsoleEndpoint}}
                  </td></tr>
                </table>
                <p style="margin:16px 0 0 0;font-size:14px;color:#4b5563;">Bucketeer Team.</p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`
)

// TemplateRenderer handles email template rendering with variable substitution
type TemplateRenderer struct {
	config Config
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(config Config) *TemplateRenderer {
	return &TemplateRenderer{config: config}
}

// getTemplateForLanguage returns the template configuration for the specified language
// Falls back to English if the requested language is not available
func (r *TemplateRenderer) getTemplateForLanguage(language string) TemplatesByLanguage {
	switch language {
	case "ja":
		return r.config.Templates.Ja
	case "en":
		fallthrough
	default:
		return r.config.Templates.En
	}
}

// RenderWelcomeEmail renders the welcome email template
func (r *TemplateRenderer) RenderWelcomeEmail(
	language string,
	userEmail string,
) (subject, body string) {
	template := r.getTemplateForLanguage(language).Welcome
	if template.Subject == "" {
		template.Subject = defaultWelcomeSubject
	}
	if template.Body == "" {
		template.Body = defaultWelcomeBody
	}

	variables := map[string]string{
		"{{baseURL}}":            r.config.BaseURL,
		"{{webConsoleEndpoint}}": r.config.BaseURL,
		"{{userEmail}}":          userEmail,
	}

	return template.Subject, r.substituteVariables(template.Body, variables)
}

// substituteVariables replaces template variables with actual values
func (r *TemplateRenderer) substituteVariables(template string, variables map[string]string) string {
	result := template
	for placeholder, value := range variables {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
