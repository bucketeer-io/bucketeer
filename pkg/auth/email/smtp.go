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
	"context"
	"fmt"
	"net/smtp"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
)

// SMTPEmailService implements EmailService using SMTP
type SMTPEmailService struct {
	config auth.EmailServiceConfig
	logger *zap.Logger
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(config auth.EmailServiceConfig, logger *zap.Logger) EmailService {
	return &SMTPEmailService{
		config: config,
		logger: logger,
	}
}

func (s *SMTPEmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
	subject := "Reset Your Bucketeer Password"
	body := s.renderPasswordResetTemplate(resetURL, resetToken)

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send password reset email",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	s.logger.Info("Password reset email sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *SMTPEmailService) SendPasswordChangedNotification(ctx context.Context, to string) error {
	subject := "Password Changed Successfully"
	body := s.renderPasswordChangedTemplate()

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send password changed notification",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send password changed notification: %w", err)
	}

	s.logger.Info("Password changed notification sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *SMTPEmailService) SendWelcomeEmail(ctx context.Context, to, tempPassword string) error {
	subject := "Welcome to Bucketeer"
	body := s.renderWelcomeTemplate(tempPassword)

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send welcome email",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	s.logger.Info("Welcome email sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *SMTPEmailService) sendEmail(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, s.config.FromEmail, subject, body))

	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
}

func (s *SMTPEmailService) renderPasswordResetTemplate(resetURL, token string) string {
	return fmt.Sprintf(`
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
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; }
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
        
        <p>We received a request to reset your Bucketeer password. If you made this request, click the button below to reset your password:</p>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="%s" class="button">Reset Password</a>
        </p>
        
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #f8f9fa; padding: 10px; border-radius: 3px;">%s</p>
        
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
</html>`, resetURL, resetURL)
}

func (s *SMTPEmailService) renderPasswordChangedTemplate() string {
	return `
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
}

func (s *SMTPEmailService) renderWelcomeTemplate(tempPassword string) string {
	return fmt.Sprintf(`
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
        <div class="temp-password">%s</div>
        
        <div class="warning">
            <strong>Important:</strong> Please change this temporary password immediately after your first login for security reasons.
        </div>
        
        <p>You can sign in at: %s</p>
        
        <p>Welcome to the team!</p>
    </div>
</body>
</html>`, tempPassword, s.config.BaseURL)
}
