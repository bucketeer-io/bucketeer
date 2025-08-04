# Summary

We will implement password authentication for Bucketeer to provide a complete authentication solution alongside the existing Google OAuth integration. This RFC proposes implementing real password-based authentication with secure password storage, password recovery via email, and proper user management.

## Background

The current authentication system in Bucketeer supports:
- **Google OAuth**: OAuth2 integration for Google authentication

However, for production deployments, especially in self-hosted environments, organizations need a native password authentication system that doesn't rely on external OAuth providers.

## Goals

- Implement secure password authentication with proper password hashing
- Add password recovery functionality via email
- Extend the existing `SignIn` API to support real password authentication
- Maintain backward compatibility with existing OAuth authentication
- Ensure security best practices for password storage and validation
- Support password complexity requirements

## Implementation

### Database Schema Changes

We need to extend the existing account system to support password authentication. Since we already have the `account_v2` table, we'll add a new table specifically for password credentials:

```sql
-- Create "account_credentials" table
CREATE TABLE `account_credentials` (
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `password_reset_token` varchar(255) DEFAULT NULL,
  `password_reset_token_expires_at` bigint DEFAULT NULL,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`email`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- Create index for password reset tokens
CREATE INDEX `idx_password_reset_token` ON `account_credentials` (`password_reset_token`);
```

### Configuration Changes

Extend the existing `OAuthConfig` structure to include password authentication settings:

```go
type PasswordAuthConfig struct {
    Enabled                    bool          `json:"enabled"`
    PasswordMinLength          int           `json:"passwordMinLength"`
    PasswordRequireUppercase   bool          `json:"passwordRequireUppercase"`
    PasswordRequireLowercase   bool          `json:"passwordRequireLowercase"`
    PasswordRequireNumbers     bool          `json:"passwordRequireNumbers"`
    PasswordRequireSymbols     bool          `json:"passwordRequireSymbols"`
    PasswordResetTokenTTL      time.Duration `json:"passwordResetTokenTTL"`
    EmailServiceEnabled        bool          `json:"emailServiceEnabled"`
    EmailServiceProvider       string        `json:"emailServiceProvider"` // "sendgrid", "smtp", etc.
    EmailServiceConfig         interface{}   `json:"emailServiceConfig"`
}

type OAuthConfig struct {
    Issuer           string               `json:"issuer"`
    Audience         string               `json:"audience"`
    GoogleConfig     GoogleConfig         `json:"google"`
    PasswordAuth     PasswordAuthConfig   `json:"passwordAuth"`
}
```

### API Extensions

#### New API Endpoints

We'll add new API endpoints for password management:

```proto
// proto/auth/service.proto

message CreatePasswordRequest {
  string email = 1;
  string password = 2;
}

message CreatePasswordResponse {}

message UpdatePasswordRequest {
  string current_password = 1;
  string new_password = 2;
}

message UpdatePasswordResponse {}

message InitiatePasswordResetRequest {
  string email = 1;
}

message InitiatePasswordResetResponse {
  string message = 1; // Generic success message for security
}

message ResetPasswordRequest {
  string reset_token = 1;
  string new_password = 2;
}

message ResetPasswordResponse {}

message ValidatePasswordResetTokenRequest {
  string reset_token = 1;
}

message ValidatePasswordResetTokenResponse {
  bool is_valid = 1;
  string email = 2;
}

service AuthService {
  // ... existing methods

  rpc CreatePassword(CreatePasswordRequest) returns (CreatePasswordResponse) {
    option (google.api.http) = {
      post: "/v1/auth/password"
      body: "*"
    };
  }

  rpc UpdatePassword(UpdatePasswordRequest) returns (UpdatePasswordResponse) {
    option (google.api.http) = {
      put: "/v1/auth/password"
      body: "*"
    };
  }

  rpc InitiatePasswordReset(InitiatePasswordResetRequest) returns (InitiatePasswordResetResponse) {
    option (google.api.http) = {
      post: "/v1/auth/password/reset/initiate"
      body: "*"
    };
  }

  rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse) {
    option (google.api.http) = {
      post: "/v1/auth/password/reset"
      body: "*"
    };
  }

  rpc ValidatePasswordResetToken(ValidatePasswordResetTokenRequest) returns (ValidatePasswordResetTokenResponse) {
    option (google.api.http) = {
      post: "/v1/auth/password/reset/validate"
      body: "*"
    };
  }
}
```

#### Enhanced SignIn Implementation

Extend the existing `SignIn` method to support both demo and password authentication:

```go
func (s *authService) SignIn(
    ctx context.Context,
    request *authproto.SignInRequest,
) (*authproto.SignInResponse, error) {
    localizer := locale.NewLocalizer(ctx)
    err := validateSignInRequest(request, localizer)
    if err != nil {
        return nil, err
    }

    // Try password authentication if enabled
    if s.config.PasswordAuth.Enabled {
        return s.handlePasswordSignIn(ctx, request, localizer)
    }

    // If password authentication is not enabled or credentials don't match, deny access
    s.logger.Error("Sign in failed - no valid authentication method",
        zap.String("email", request.Email),
    )
    dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
        Locale:  localizer.GetLocale(),
        Message: localizer.MustLocalize(locale.PermissionDenied),
    })
    if err != nil {
        return nil, err
    }
    return nil, dt.Err()
}

func (s *authService) handlePasswordSignIn(
    ctx context.Context,
    request *authproto.SignInRequest,
    localizer locale.Localizer,
) (*authproto.SignInResponse, error) {
    // Verify password
    if !s.verifyPassword(ctx, request.Email, request.Password) {
        s.logger.Error("Password sign in failed - invalid credentials",
            zap.String("email", request.Email),
        )
        dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
            Locale:  localizer.GetLocale(),
            Message: localizer.MustLocalize(locale.PermissionDenied),
        })
        if err != nil {
            return nil, err
        }
        return nil, dt.Err()
    }

    // Get organizations for the user
    organizations, err := s.getOrganizationsByEmail(ctx, request.Email, localizer)
    if err != nil {
        return nil, err
    }

    // Check account status - user should have at least one enabled account
    account, err := s.checkAccountStatus(ctx, request.Email, organizations, localizer)
    if err != nil {
        s.logger.Error("Failed to check account status",
            zap.Error(err),
            zap.String("email", request.Email),
        )
        return nil, err
    }

    accountDomain := domain.AccountV2{AccountV2: account.Account}
    isSystemAdmin := s.hasSystemAdminOrganization(organizations)

    // Generate token for successful authentication
    token, err := s.generateToken(ctx, request.Email, accountDomain, isSystemAdmin, localizer)
    if err != nil {
        return nil, err
    }

    s.logger.Info("Successful password authentication",
        zap.String("email", request.Email),
    )

    return &authproto.SignInResponse{Token: token}, nil
}
```

### Password Security Implementation

#### Password Hashing

We'll use bcrypt for password hashing, which is considered secure and includes automatic salting:

```go
package auth

import (
    "golang.org/x/crypto/bcrypt"
)

const (
    // bcrypt cost factor - balance between security and performance
    BcryptCost = 14
)

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func ValidatePassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### Password Validation

```go
func ValidatePasswordComplexity(password string, config PasswordAuthConfig) error {
    if len(password) < config.PasswordMinLength {
        return errors.New("password too short")
    }
    
    if config.PasswordRequireUppercase && !containsUppercase(password) {
        return errors.New("password must contain uppercase letters")
    }
    
    if config.PasswordRequireLowercase && !containsLowercase(password) {
        return errors.New("password must contain lowercase letters")
    }
    
    if config.PasswordRequireNumbers && !containsNumbers(password) {
        return errors.New("password must contain numbers")
    }
    
    if config.PasswordRequireSymbols && !containsSymbols(password) {
        return errors.New("password must contain symbols")
    }
    
    return nil
}
```



### Email Service Integration

#### Email Service Interface

```go
package auth

type EmailService interface {
    SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error
    SendPasswordChangedNotification(ctx context.Context, to string) error
    SendWelcomeEmail(ctx context.Context, to, tempPassword string) error
}

type EmailServiceConfig struct {
    Provider         string `json:"provider"`         // "smtp", "sendgrid", "ses"
    SMTPHost         string `json:"smtpHost"`
    SMTPPort         int    `json:"smtpPort"`
    SMTPUsername     string `json:"smtpUsername"`
    SMTPPassword     string `json:"smtpPassword"`
    SendGridAPIKey   string `json:"sendgridAPIKey"`
    SESRegion        string `json:"sesRegion"`
    SESAccessKey     string `json:"sesAccessKey"`
    SESSecretKey     string `json:"sesSecretKey"`
    FromEmail        string `json:"fromEmail"`
    FromName         string `json:"fromName"`
    BaseURL          string `json:"baseURL"`          // For constructing reset URLs
}
```

#### Email Service Providers

We'll support multiple email service providers to accommodate different deployment scenarios:

##### 1. SMTP (Simple Mail Transfer Protocol)
**Best for**: Self-hosted deployments, corporate environments with existing mail servers

```go
type SMTPEmailService struct {
    config EmailServiceConfig
    logger *zap.Logger
}

func (s *SMTPEmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
    subject := "Reset Your Bucketeer Password"
    body := s.renderPasswordResetTemplate(resetURL, resetToken)
    
    return s.sendEmail(ctx, to, subject, body)
}

func (s *SMTPEmailService) sendEmail(ctx context.Context, to, subject, body string) error {
    auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
    
    msg := []byte(fmt.Sprintf("To: %s\r\n"+
        "Subject: %s\r\n"+
        "Content-Type: text/html; charset=UTF-8\r\n"+
        "\r\n"+
        "%s\r\n", to, subject, body))
    
    addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
    return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
}
```

##### 2. SendGrid
**Best for**: Cloud deployments, scalable email delivery, advanced analytics

```go
type SendGridEmailService struct {
    client *sendgrid.Client
    config EmailServiceConfig
    logger *zap.Logger
}

func NewSendGridEmailService(config EmailServiceConfig, logger *zap.Logger) *SendGridEmailService {
    return &SendGridEmailService{
        client: sendgrid.NewSendClient(config.SendGridAPIKey),
        config: config,
        logger: logger,
    }
}

func (s *SendGridEmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
    from := mail.NewEmail(s.config.FromName, s.config.FromEmail)
    toEmail := mail.NewEmail("", to)
    subject := "Reset Your Bucketeer Password"
    
    plainTextContent := fmt.Sprintf("Reset your password by clicking: %s", resetURL)
    htmlContent := s.renderPasswordResetTemplate(resetURL, resetToken)
    
    message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
    
    response, err := s.client.Send(message)
    if err != nil {
        s.logger.Error("Failed to send email via SendGrid",
            zap.Error(err),
            zap.String("to", to),
        )
        return err
    }
    
    if response.StatusCode >= 400 {
        s.logger.Error("SendGrid returned error status",
            zap.Int("statusCode", response.StatusCode),
            zap.String("body", response.Body),
        )
        return fmt.Errorf("sendgrid error: %d", response.StatusCode)
    }
    
    return nil
}
```

##### 3. Amazon SES
**Best for**: AWS deployments, cost-effective bulk email sending

```go
type SESEmailService struct {
    client *ses.Client
    config EmailServiceConfig
    logger *zap.Logger
}

func NewSESEmailService(config EmailServiceConfig, logger *zap.Logger) (*SESEmailService, error) {
    cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
        awsconfig.WithRegion(config.SESRegion),
        awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            config.SESAccessKey,
            config.SESSecretKey,
            "",
        )),
    )
    if err != nil {
        return nil, err
    }
    
    return &SESEmailService{
        client: ses.NewFromConfig(cfg),
        config: config,
        logger: logger,
    }, nil
}

func (s *SESEmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
    subject := "Reset Your Bucketeer Password"
    htmlBody := s.renderPasswordResetTemplate(resetURL, resetToken)
    textBody := fmt.Sprintf("Reset your password by visiting: %s", resetURL)
    
    input := &ses.SendEmailInput{
        Destination: &sestypes.Destination{
            ToAddresses: []string{to},
        },
        Message: &sestypes.Message{
            Body: &sestypes.Body{
                Html: &sestypes.Content{
                    Charset: aws.String("UTF-8"),
                    Data:    aws.String(htmlBody),
                },
                Text: &sestypes.Content{
                    Charset: aws.String("UTF-8"),
                    Data:    aws.String(textBody),
                },
            },
            Subject: &sestypes.Content{
                Charset: aws.String("UTF-8"),
                Data:    aws.String(subject),
            },
        },
        Source: aws.String(s.config.FromEmail),
    }
    
    _, err := s.client.SendEmail(ctx, input)
    if err != nil {
        s.logger.Error("Failed to send email via SES",
            zap.Error(err),
            zap.String("to", to),
        )
        return err
    }
    
    return nil
}
```

#### Email Templates

We'll create professional, security-focused email templates:

```go
func (s *EmailServiceBase) renderPasswordResetTemplate(resetURL, token string) string {
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

func (s *EmailServiceBase) renderPasswordChangedTemplate() string {
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

func (s *EmailServiceBase) renderWelcomeTemplate(tempPassword string) string {
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
        
        <p>You can sign in at: [Your Bucketeer URL]</p>
        
        <p>Welcome to the team!</p>
    </div>
</body>
</html>`, tempPassword)
}
```

#### Email Service Factory

```go
func NewEmailService(config EmailServiceConfig, logger *zap.Logger) (EmailService, error) {
    switch config.Provider {
    case "smtp":
        return NewSMTPEmailService(config, logger), nil
    case "sendgrid":
        return NewSendGridEmailService(config, logger), nil
    case "ses":
        return NewSESEmailService(config, logger)
    default:
        return nil, fmt.Errorf("unsupported email provider: %s", config.Provider)
    }
}
```

#### Password Reset Implementation

```go
func (s *authService) InitiatePasswordReset(
    ctx context.Context,
    request *authproto.InitiatePasswordResetRequest,
) (*authproto.InitiatePasswordResetResponse, error) {
    localizer := locale.NewLocalizer(ctx)
    
    // Validate that the user has organizations (i.e., account exists)
    organizations, err := s.getOrganizationsByEmail(ctx, request.Email, localizer)
    if err != nil {
        // For security, don't reveal whether the account exists
        s.logger.Warn("Password reset attempted for non-existent account",
            zap.String("email", request.Email),
        )
        return &authproto.InitiatePasswordResetResponse{
            Message: "If an account with this email exists, a password reset link has been sent.",
        }, nil
    }
    
    // Generate secure reset token
    resetToken, err := s.generateSecureToken()
    if err != nil {
        return nil, auth.StatusInternal
    }
    
    // Store reset token with expiration
    expiresAt := time.Now().Add(s.config.PasswordAuth.PasswordResetTokenTTL).Unix()
    err = s.credentialsStorage.SetPasswordResetToken(
        ctx, 
        request.Email, 
        resetToken, 
        expiresAt,
    )
    if err != nil {
        return nil, auth.StatusInternal
    }
    
    // Send reset email
    resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.config.BaseURL, resetToken)
    err = s.emailService.SendPasswordResetEmail(ctx, request.Email, resetToken, resetURL)
    if err != nil {
        s.logger.Error("Failed to send password reset email",
            zap.Error(err),
            zap.String("email", request.Email),
        )
        // Don't return error to user for security reasons
    }
    
    return &authproto.InitiatePasswordResetResponse{
        Message: "If an account with this email exists, a password reset link has been sent.",
    }, nil
}
```

### Storage Layer

#### Credentials Storage Interface

```go
package storage

type CredentialsStorage interface {
    CreateCredentials(ctx context.Context, email, passwordHash string) error
    GetCredentials(ctx context.Context, email string) (*domain.AccountCredentials, error)
    UpdatePassword(ctx context.Context, email, passwordHash string) error
    DeleteCredentials(ctx context.Context, email string) error
    
    SetPasswordResetToken(ctx context.Context, email, token string, expiresAt int64) error
    GetPasswordResetToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
    DeletePasswordResetToken(ctx context.Context, token string) error
}
```

#### Domain Models

```go
package domain

type AccountCredentials struct {
    Email                      string
    PasswordHash              string
    CreatedAt                 int64
    UpdatedAt                 int64
}

type PasswordResetToken struct {
    Token          string
    Email          string
    ExpiresAt      int64
    CreatedAt      int64
}
```

### Error Handling

Add new error codes to the auth package:

```go
// pkg/auth/error.go

var (
    // ... existing errors
    
    // Password-related errors
    StatusPasswordTooWeak         = gstatus.New(codes.InvalidArgument, "auth: password too weak")
    StatusPasswordMismatch        = gstatus.New(codes.InvalidArgument, "auth: password mismatch")
    StatusPasswordAlreadyExists   = gstatus.New(codes.AlreadyExists, "auth: password already exists")
    StatusPasswordNotFound        = gstatus.New(codes.NotFound, "auth: password not found")
    
    // Password reset errors
    StatusInvalidResetToken       = gstatus.New(codes.InvalidArgument, "auth: invalid reset token")
    StatusExpiredResetToken       = gstatus.New(codes.InvalidArgument, "auth: reset token expired")
    StatusResetTokenNotFound      = gstatus.New(codes.NotFound, "auth: reset token not found")
    
    // Email service errors
    StatusEmailServiceUnavailable = gstatus.New(codes.Unavailable, "auth: email service unavailable")
    StatusTooManyEmailRequests    = gstatus.New(codes.ResourceExhausted, "auth: too many email requests")
    StatusInvalidEmailConfig      = gstatus.New(codes.InvalidArgument, "auth: invalid email configuration")
)
```

## Security Considerations

### Password Storage
- Use bcrypt with a cost factor of 12 for password hashing
- Never store passwords in plain text
- Use secure random salt generation (bcrypt handles this automatically)

### Rate Limiting
- Log all authentication attempts for monitoring
- Consider implementing rate limiting at the application/network level if needed

### Token Security
- Use cryptographically secure random token generation
- Set reasonable expiration times for reset tokens (e.g., 1 hour)
- Invalidate tokens after use
- Store tokens securely (hashed if possible)

### Input Validation
- Validate all input parameters
- Implement password complexity requirements
- Sanitize email addresses
- Prevent timing attacks in authentication checks

## Configuration Examples

### SMTP Configuration (Self-hosted)
**Recommended for**: Corporate environments, self-hosted deployments

```yaml
# values.yaml
auth:
  passwordAuth:
    enabled: true
    passwordMinLength: 8
    passwordRequireUppercase: true
    passwordRequireLowercase: true
    passwordRequireNumbers: true
    passwordRequireSymbols: false
    passwordResetTokenTTL: "1h"
    emailServiceEnabled: true
    emailServiceConfig:
      provider: "smtp"
      smtpHost: "smtp.example.com"
      smtpPort: 587
      smtpUsername: "noreply@example.com"
      smtpPassword: "${SMTP_PASSWORD}"
      fromEmail: "noreply@example.com"
      fromName: "Bucketeer"
      baseURL: "https://bucketeer.example.com"

# Gmail SMTP Example
auth:
  passwordAuth:
    emailServiceConfig:
      provider: "smtp"
      smtpHost: "smtp.gmail.com"
      smtpPort: 587
      smtpUsername: "your-email@gmail.com"
      smtpPassword: "${GMAIL_APP_PASSWORD}"  # Use App Password, not regular password
      fromEmail: "your-email@gmail.com"
      fromName: "Bucketeer"
```

### SendGrid Configuration (Cloud)
**Recommended for**: Production cloud deployments, high-volume email

```yaml
# values.yaml
auth:
  passwordAuth:
    enabled: true
    passwordMinLength: 8
    passwordRequireUppercase: true
    passwordRequireLowercase: true
    passwordRequireNumbers: true
    passwordRequireSymbols: false
    passwordResetTokenTTL: "1h"
    emailServiceEnabled: true
    emailServiceConfig:
      provider: "sendgrid"
      sendgridAPIKey: "${SENDGRID_API_KEY}"
      fromEmail: "noreply@yourdomain.com"
      fromName: "Bucketeer"
      baseURL: "https://bucketeer.yourdomain.com"
```

### Amazon SES Configuration (AWS)
**Recommended for**: AWS deployments, cost-effective scaling

```yaml
# values.yaml
auth:
  passwordAuth:
    enabled: true
    passwordMinLength: 8
    passwordRequireUppercase: true
    passwordRequireLowercase: true
    passwordRequireNumbers: true
    passwordRequireSymbols: false
    passwordResetTokenTTL: "1h"
    emailServiceEnabled: true
    emailServiceConfig:
      provider: "ses"
      sesRegion: "us-east-1"
      sesAccessKey: "${AWS_ACCESS_KEY_ID}"
      sesSecretKey: "${AWS_SECRET_ACCESS_KEY}"
      fromEmail: "noreply@yourdomain.com"
      fromName: "Bucketeer"
      baseURL: "https://bucketeer.yourdomain.com"
```

### Email Service Provider Comparison

| Provider | Cost | Setup Complexity | Reliability | Features | Best For |
|----------|------|------------------|-------------|----------|----------|
| **SMTP** | Low (if you have mail server) | Medium | Good | Basic | Corporate/Self-hosted |
| **SendGrid** | $$$ (pay per email) | Low | Excellent | Advanced analytics, templates | Production cloud |
| **Amazon SES** | $ (very cheap) | Medium | Excellent | Basic, AWS integration | AWS deployments |

### Recommended Email Providers by Deployment Type

#### Self-Hosted Deployments
1. **Corporate SMTP Server** - If available
2. **Gmail SMTP** - For small teams (use App Passwords)
3. **Amazon SES** - Cost-effective for any scale

#### Cloud Deployments
1. **SendGrid** - Best for production with analytics
2. **Amazon SES** - Most cost-effective
3. **Cloud provider email services** (Google Cloud, Azure)

#### Development/Testing
1. **Mailhog** or **MailCatcher** - Local email testing
2. **Gmail SMTP** - Simple setup for development
3. **SendGrid Free Tier** - 100 emails/day free

### Email Security and Rate Limiting

#### Email Rate Limiting
To prevent abuse, implement rate limiting for password reset emails:

```go
type EmailRateLimiter struct {
    redis       redis.Client
    maxAttempts int
    window      time.Duration
}

func (e *EmailRateLimiter) CheckRateLimit(ctx context.Context, email string) error {
    key := fmt.Sprintf("email_rate_limit:%s", email)
    
    count, err := e.redis.Get(ctx, key).Int()
    if err != nil && err != redis.Nil {
        return err
    }
    
    if count >= e.maxAttempts {
        return auth.StatusTooManyEmailRequests
    }
    
    // Increment counter
    pipe := e.redis.Pipeline()
    pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, e.window)
    _, err = pipe.Exec(ctx)
    
    return err
}
```

#### Email Template Security

```go
func sanitizeEmailContent(content string) string {
    // Remove any potential XSS vectors
    content = html.EscapeString(content)
    // Remove script tags if any
    re := regexp.MustCompile(`<script[^>]*>.*?</script>`)
    content = re.ReplaceAllString(content, "")
    return content
}
```

## Future Enhancements

### Multi-Factor Authentication (MFA)
- TOTP (Time-based One-Time Password) support
- SMS verification
- Hardware token support

### Advanced Security Features
- Password history to prevent reuse
- Forced password rotation policies
- Advanced anomaly detection
- Single Sign-On (SSO) integration with other providers

### User Experience Improvements
- Password strength indicators
- Social login integration
- Remember me functionality
- Account recovery options

## Conclusion

This implementation provides a secure, production-ready password authentication system for Bucketeer while maintaining compatibility with existing authentication methods. The phased approach ensures minimal disruption to current users while providing a solid foundation for future security enhancements. 