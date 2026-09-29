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

// Template represents a single email template with subject and body
type Template struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// TemplatesByLanguage represents email templates for a specific language
type TemplatesByLanguage struct {
	PasswordChanged Template `json:"passwordChanged"`
	PasswordSetup   Template `json:"passwordSetup"`
	PasswordReset   Template `json:"passwordReset"`
	Welcome         Template `json:"welcome"`
}

// TemplatesConfig represents email templates for all supported languages
type TemplatesConfig struct {
	En TemplatesByLanguage `json:"en"`
	Ja TemplatesByLanguage `json:"ja"`
}

// SMTPConfig represents SMTP server configuration
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// SendGridConfig represents SendGrid API configuration
type SendGridConfig struct {
	APIKey string `json:"apiKey"`
}

// SESConfig represents Amazon SES configuration
type SESConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// MailerSendConfig represents MailerSend API configuration
type MailerSendConfig struct {
	APIKey string `json:"apiKey"`
}

// SenderConfig represents the sender's email and name
type SenderConfig struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Config represents the email service configuration
type Config struct {
	Enabled    bool             `json:"enabled"`
	Provider   string           `json:"provider"` // "smtp", "sendgrid", "ses", "mailersend"
	SMTP       SMTPConfig       `json:"smtp"`
	SendGrid   SendGridConfig   `json:"sendgrid"`
	SES        SESConfig        `json:"ses"`
	MailerSend MailerSendConfig `json:"mailersend"`
	Sender     SenderConfig     `json:"sender"`
	BaseURL    string           `json:"baseURL"` // Base URL for web console
	Templates  TemplatesConfig  `json:"templates"`
}
