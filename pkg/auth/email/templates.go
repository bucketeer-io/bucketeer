// Copyright 2025 The Bucketeer Authors.
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
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
)

// Default email templates
var (
	defaultPasswordChangedSubject = "Your Bucketeer Password Was Changed"
	//nolint:lll
	defaultPasswordChangedBody = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="x-apple-disable-message-reformatting">
        <title>Password Changed</title>
    </head>
    <body style="margin:0;padding:0;background:#f5f5f5;color:#333;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;line-height:1.6;">
        <div style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0;">
            If you didn't make this change, contact support immediately.
        </div>
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
            <tr>
                <td align="center" style="padding:24px 12px;">
                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="width:100%;max-width:600px;background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
                        <tr>
                            <td style="padding:32px;background:#ffffff;">
                                <img src="{{webConsoleEndpoint}}/img/bucketeer-logo-primary.png" alt="Bucketeer" width="205" height="48" style="display:block;margin-bottom:24px;" />
                                <p style="margin:0 0 24px 0;font-size:15px;color:#4b5563;">Hello, {{userEmail}}.</p>
                                <p style="margin:0 0 24px 0;font-size:15px;color:#4b5563;">We wanted to let you know that your Bucketeer password has changed.</p>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td style="padding:16px 20px;background:#FDFBE8;border-left:4px solid #FFB802;border-radius:6px;">
                                            <p style="margin:0;font-size:14px;color:#725201;line-height:1.5;"><strong>Security notice:</strong> If you didn't make this change, please contact your administrator immediately.</p>
                                        </td>
                                    </tr>
                                </table>
                                <p style="margin:0 0 12px 0;font-size:15px;font-weight:600;color:#4b5563;">For your security:</p>
                                <ul style="margin:0 0 24px 0;padding-left:20px;font-size:14px;color:#4b5563;line-height:1.8;">
                                    <li style="margin-bottom:8px;">Use a strong, unique password</li>
                                    <li style="margin-bottom:8px;">Never share your password with anyone</li>
                                    <li>Consider using a password manager</li>
                                </ul>
                                <p style="margin:0;font-size:14px;color:#4b5563;">Thank you for keeping your account secure.</p>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </body>
</html>`

	defaultPasswordSetupSubject = "Complete Your Bucketeer Password Setup"
	//nolint:lll
	defaultPasswordSetupBody = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="x-apple-disable-message-reformatting">
        <title>Complete Password Setup</title>
    </head>
    <body style="margin:0;padding:0;background:#f5f5f5;color:#333;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;line-height:1.6;">
        <div style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0;">
            Set up your Bucketeer password to get started.
        </div>
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
            <tr>
                <td align="center" style="padding:24px 12px;">
                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="width:100%;max-width:600px;background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
                        <tr>
                            <td style="padding:32px;background:#ffffff;">
                                <img src="{{webConsoleEndpoint}}/img/bucketeer-logo-primary.png" alt="Bucketeer" width="205" height="48" style="display:block;margin-bottom:24px;" />
                                <p style="margin:0 0 24px 0;font-size:15px;color:#4b5563;">Hello, {{userEmail}}.</p>
                                <p style="margin:0 0 24px 0;font-size:15px;color:#4b5563;">Your Bucketeer account is ready! To complete your account setup, please create a password by clicking the button below:</p>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td align="center">
                                            <a href="{{setupURL}}" style="background:#5B21B6;color:#ffffff;text-decoration:none;padding:14px 32px;border-radius:6px;display:inline-block;font-weight:600;font-size:15px;">Set Up Password</a>
                                        </td>
                                    </tr>
                                </table>
                                <p style="margin:0 0 8px 0;font-size:13px;color:#4b5563;">Or copy and paste this link into your browser:</p>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td style="padding:12px 16px;background:#f9fafb;border:1px solid #e5e7eb;border-radius:6px;word-break:break-all;font-size:13px;color:#4b5563;">
                                            {{setupURL}}
                                        </td>
                                    </tr>
                                </table>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td style="padding:16px 20px;background:#ECF6FD;border-left:4px solid #399CE4;border-radius:6px;">
                                            <p style="margin:0 0 12px 0;font-size:14px;font-weight:600;color:#23405D;">Important:</p>
                                            <ul style="margin:0;padding-left:20px;font-size:14px;color:#29577F;line-height:1.8;">
                                                <li style="margin-bottom:8px;">This link will expire in {{expirationTime}} for security reasons</li>
                                                <li style="margin-bottom:8px;">Setting up a password will allow you to sign in directly without OAuth</li>
                                                <li>You can continue using Google sign-in even after setting up a password</li>
                                            </ul>
                                        </td>
                                    </tr>
                                </table>
                                <p style="margin:0;font-size:14px;color:#4b5563;">If you have any questions, please contact your system administrator.</p>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </body>
</html>`

	defaultPasswordResetSubject = "Reset Your Bucketeer Password"
	//nolint:lll
	defaultPasswordResetBody = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="x-apple-disable-message-reformatting">
        <title>Reset Password</title>
    </head>
    <body style="margin:0;padding:0;background:#f5f5f5;color:#333;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;line-height:1.6;">
        <div style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0;">
            Reset your password to regain access to your account.
        </div>
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
            <tr>
                <td align="center" style="padding:24px 12px;">
                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="width:100%;max-width:600px;background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
                        <tr>
                            <td style="padding:32px;background:#ffffff;">
                                <img src="{{webConsoleEndpoint}}/img/bucketeer-logo-primary.png" alt="Bucketeer" width="205" height="48" style="display:block;margin-bottom:24px;" />
                                <p style="margin:0 0 24px 0;font-size:15px;color:#4b5563;">Hello, {{userEmail}}.</p>
                                <p style="margin:0 0 24px 0;font-size:15px;color:#4b5563;">We received a request to reset your Bucketeer password. Click the button below to create a new password:</p>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td align="center">
                                            <a href="{{resetURL}}" style="background:#5B21B6;color:#ffffff;text-decoration:none;padding:14px 32px;border-radius:6px;display:inline-block;font-weight:600;font-size:15px;">Reset Password</a>
                                        </td>
                                    </tr>
                                </table>
                                <p style="margin:0 0 8px 0;font-size:13px;color:#4b5563;">Or copy and paste this link into your browser:</p>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td style="padding:12px 16px;background:#f9fafb;border:1px solid #e5e7eb;border-radius:6px;word-break:break-all;font-size:13px;color:#4b5563;">
                                            {{resetURL}}
                                        </td>
                                    </tr>
                                </table>
                                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:0 0 24px 0;">
                                    <tr>
                                        <td style="padding:16px 20px;background:#ECF6FD;border-left:4px solid #399CE4;border-radius:6px;">
                                            <p style="margin:0 0 12px 0;font-size:14px;font-weight:600;color:#23405D;">Security Information:</p>
                                            <ul style="margin:0;padding-left:20px;font-size:14px;color:#29577F;line-height:1.8;">
                                                <li style="margin-bottom:8px;">This link will expire in {{expirationTime}} for security reasons</li>
                                                <li style="margin-bottom:8px;">If you did not request this password reset, please ignore this email</li>
                                                <li>Your password will remain unchanged until you create a new one</li>
                                            </ul>
                                        </td>
                                    </tr>
                                </table>
                                <p style="margin:0;font-size:14px;color:#4b5563;">If you have any questions, please contact your system administrator.</p>
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
	config auth.EmailServiceConfig
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(config auth.EmailServiceConfig) *TemplateRenderer {
	return &TemplateRenderer{config: config}
}

// getTemplateForLanguage returns the template configuration for the specified language
// Falls back to English if the requested language is not available
func (r *TemplateRenderer) getTemplateForLanguage(language string) auth.EmailTemplatesByLanguage {
	switch language {
	case "ja":
		return r.config.Templates.Ja
	case "en":
		fallthrough
	default:
		return r.config.Templates.En
	}
}

// RenderPasswordChangedEmail renders the password changed notification template
func (r *TemplateRenderer) RenderPasswordChangedEmail(
	language string,
	userEmail string,
) (subject, body string) {
	template := r.getTemplateForLanguage(language).PasswordChanged
	if template.Subject == "" {
		template.Subject = defaultPasswordChangedSubject
	}
	if template.Body == "" {
		template.Body = defaultPasswordChangedBody
	}

	variables := map[string]string{
		"{{baseURL}}":            r.config.BaseURL,
		"{{webConsoleEndpoint}}": r.config.BaseURL,
		"{{userEmail}}":          userEmail,
	}

	return template.Subject, r.substituteVariables(template.Body, variables)
}

// RenderPasswordSetupEmail renders the password setup email template
func (r *TemplateRenderer) RenderPasswordSetupEmail(
	language string,
	setupURL string,
	ttl time.Duration,
	userEmail string,
) (subject, body string) {
	template := r.getTemplateForLanguage(language).PasswordSetup
	if template.Subject == "" {
		template.Subject = defaultPasswordSetupSubject
	}
	if template.Body == "" {
		template.Body = defaultPasswordSetupBody
	}

	variables := map[string]string{
		"{{setupURL}}":           setupURL,
		"{{baseURL}}":            r.config.BaseURL,
		"{{webConsoleEndpoint}}": r.config.BaseURL,
		"{{expirationTime}}":     ttl.String(),
		"{{userEmail}}":          userEmail,
	}

	return template.Subject, r.substituteVariables(template.Body, variables)
}

// RenderPasswordResetEmail renders the password reset email template
func (r *TemplateRenderer) RenderPasswordResetEmail(
	language string,
	resetURL string,
	ttl time.Duration,
	userEmail string,
) (subject, body string) {
	template := r.getTemplateForLanguage(language).PasswordReset
	if template.Subject == "" {
		template.Subject = defaultPasswordResetSubject
	}
	if template.Body == "" {
		template.Body = defaultPasswordResetBody
	}

	variables := map[string]string{
		"{{resetURL}}":           resetURL,
		"{{baseURL}}":            r.config.BaseURL,
		"{{webConsoleEndpoint}}": r.config.BaseURL,
		"{{expirationTime}}":     ttl.String(),
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
