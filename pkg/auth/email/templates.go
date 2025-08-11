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

	"github.com/bucketeer-io/bucketeer/pkg/auth"
)

// Default email templates
var (
	defaultPasswordResetSubject = "Reset Your Bucketeer Password"
	defaultPasswordResetBody    = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Bucketeer Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; 
                  color: white; text-decoration: none; border-radius: 5px; }
        .warning { background-color: #fff3cd; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .footer { font-size: 12px; color: #666; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Reset Your Bucketeer Password</h1>
        </div>
        
        <p>Hello,</p>
        
        <p>We received a request to reset your Bucketeer password. If you made this request, 
        click the button below to reset your password:</p>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="{{resetURL}}" class="button">Reset Password</a>
        </p>
        
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #f8f9fa; padding: 10px; border-radius: 3px;">{{resetURL}}</p>
        
        <div class="warning">
            <strong>Security Note:</strong>
            <ul>
                <li>This link will expire in 1 hour for security reasons</li>
                <li>If you didn't request this password reset, please ignore this email</li>
                <li>Never share this link with anyone</li>
            </ul>
        </div>
        
        <div class="footer">
            <p>This is an automated message from Bucketeer. Please do not reply to this email.</p>
            <p>If you have any questions, please contact your system administrator.</p>
        </div>
    </div>
</body>
</html>`

	defaultPasswordChangedSubject = "Password Changed Successfully"
	defaultPasswordChangedBody    = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Password Changed Successfully</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #d4edda; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .alert { background-color: #fff3cd; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Password Changed Successfully</h1>
        </div>
        
        <p>Hello,</p>
        
        <p>This email confirms that your Bucketeer password has been successfully changed.</p>
        
        <div class="alert">
            <strong>Security Notice:</strong>
            If you did not make this change, please contact your system administrator immediately.
        </div>
        
        <p>For your security:</p>
        <ul>
            <li>Always use a strong, unique password</li>
            <li>Never share your password with anyone</li>
            <li>Consider using a password manager</li>
        </ul>
        
        <p>Thank you for keeping your account secure.</p>
    </div>
</body>
</html>`

	defaultPasswordSetupSubject = "Complete Your Bucketeer Password Setup"
	defaultPasswordSetupBody    = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Complete Your Bucketeer Password Setup</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; 
                  color: white; text-decoration: none; border-radius: 5px; }
        .info { background-color: #d1ecf1; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .footer { font-size: 12px; color: #666; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Complete Your Bucketeer Password Setup</h1>
        </div>
        
        <p>Hello,</p>
        
        <p>Your Bucketeer account is ready! To complete your account setup, please create 
        a password by clicking the button below:</p>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="{{setupURL}}" class="button">Set Up Password</a>
        </p>
        
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #f8f9fa; padding: 10px; border-radius: 3px;">{{setupURL}}</p>
        
        <div class="info">
            <strong>Important:</strong>
            <ul>
                <li>This link will expire in 24 hours for security reasons</li>
                <li>Setting up a password will allow you to sign in directly without OAuth</li>
                <li>You can continue using Google sign-in even after setting up a password</li>
            </ul>
        </div>
        
        <div class="footer">
            <p>This is an automated message from Bucketeer. Please do not reply to this email.</p>
            <p>If you have any questions, please contact your system administrator.</p>
        </div>
    </div>
</body>
</html>`

	defaultWelcomeSubject = "Welcome to Bucketeer"
	defaultWelcomeBody    = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome to Bucketeer</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .temp-password { background-color: #f8f9fa; padding: 15px; border-radius: 5px; font-family: monospace; }
        .warning { background-color: #f8d7da; padding: 15px; border-radius: 5px; margin: 20px 0; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Bucketeer!</h1>
        </div>
        
        <p>Hello,</p>
        
        <p>Your Bucketeer account has been created. Here are your login credentials:</p>
        
        <p><strong>Temporary Password:</strong></p>
        <div class="temp-password">{{tempPassword}}</div>
        
        <div class="warning">
            <strong>Important:</strong> Please change this temporary password immediately 
            after your first login for security reasons.
        </div>
        
        <p>You can sign in at: {{baseURL}}</p>
        
        <p>Welcome to the team!</p>
    </div>
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

// RenderPasswordResetEmail renders the password reset email template
func (r *TemplateRenderer) RenderPasswordResetEmail(resetURL, resetToken string) (subject, body string) {
	template := r.config.Templates.PasswordReset
	if template.Subject == "" {
		template.Subject = defaultPasswordResetSubject
	}
	if template.Body == "" {
		template.Body = defaultPasswordResetBody
	}

	variables := map[string]string{
		"{{resetURL}}":   resetURL,
		"{{resetToken}}": resetToken,
		"{{baseURL}}":    r.config.BaseURL,
	}

	return template.Subject, r.substituteVariables(template.Body, variables)
}

// RenderPasswordChangedEmail renders the password changed notification template
func (r *TemplateRenderer) RenderPasswordChangedEmail() (subject, body string) {
	template := r.config.Templates.PasswordChanged
	if template.Subject == "" {
		template.Subject = defaultPasswordChangedSubject
	}
	if template.Body == "" {
		template.Body = defaultPasswordChangedBody
	}

	variables := map[string]string{
		"{{baseURL}}": r.config.BaseURL,
	}

	return template.Subject, r.substituteVariables(template.Body, variables)
}

// RenderPasswordSetupEmail renders the password setup email template
func (r *TemplateRenderer) RenderPasswordSetupEmail(setupURL, setupToken string) (subject, body string) {
	template := r.config.Templates.PasswordSetup
	if template.Subject == "" {
		template.Subject = defaultPasswordSetupSubject
	}
	if template.Body == "" {
		template.Body = defaultPasswordSetupBody
	}

	variables := map[string]string{
		"{{setupURL}}":   setupURL,
		"{{setupToken}}": setupToken,
		"{{baseURL}}":    r.config.BaseURL,
	}

	return template.Subject, r.substituteVariables(template.Body, variables)
}

// RenderWelcomeEmail renders the welcome email template
func (r *TemplateRenderer) RenderWelcomeEmail(tempPassword string) (subject, body string) {
	template := r.config.Templates.Welcome
	if template.Subject == "" {
		template.Subject = defaultWelcomeSubject
	}
	if template.Body == "" {
		template.Body = defaultWelcomeBody
	}

	variables := map[string]string{
		"{{tempPassword}}": tempPassword,
		"{{baseURL}}":      r.config.BaseURL,
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
