# Magic Link Authentication Flow - Design Document

## Overview

This document outlines the design for a new authentication flow that uses email verification via magic links before
organization selection. This approach improves security by preventing user enumeration attacks and ensures email
ownership verification before granting access.

## Problem Statement

### Current Flow

**OAuth (Google):**

1. User clicks "Sign in with Google"
2. User authenticates with Google → `ExchangeToken` → Token issued for first organization
3. User accesses console with that organization

**Password:**

1. User enters email + password → `SignIn` → Token issued for first organization
2. User accesses console with that organization

**Organization Switching:**

- System admins can switch organizations using `SwitchOrganization` API

### Security Issues

- **User Enumeration**: Attackers can determine if an email exists in the system by observing different error messages
  or response times during authentication attempts
- **No Organization Choice**: Users are automatically assigned to their first organization without selecting which org
  they want to access
- **Organization-Blind Authentication**: Users authenticate before knowing which org they're accessing, which may have
  different security requirements
- **Mixed Auth Methods**: No way to enforce different authentication methods per organization during login (only after
  token is issued)

## Proposed Solution

### New Authentication Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. Email Submission                                             │
│    POST /v1/auth/request_magic_link                             │
│    { "email": "user@example.com" }                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Generic Response (always 200 OK)                             │
│    "If this email is registered, we've sent instructions."      │
│    (Prevents enumeration - no indication if email exists)       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Backend Processing                                           │
│    - Create email verification token (always)                   │
│    - IF email exists: Send magic link email                     │
│    - IF email doesn't exist: Silent failure (no email sent)     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. User Clicks Magic Link                                       │
│    POST /v1/auth/verify_magic_link                              │
│    { "token": "..." }                                           │
│    Returns:                                                     │
│    - email (verified)                                           │
│    - organizations list with authentication_settings            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. Organization Selection                                       │
│    Frontend displays organizations from step 4 response         │
│    User selects organization from UI                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 6. Show Auth Methods for Selected Org                           │
│    Frontend reads authentication_settings.enabled_types:        │
│    - AUTHENTICATION_TYPE_PASSWORD (2)                           │
│    - AUTHENTICATION_TYPE_GOOGLE (1)                             │
│    Display available methods to user                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 7. Complete Authentication                                      │
│                                                                 │
│    If Password:                                                 │
│      POST /v1/auth/signin                                       │
│      { "email": "...", "password": "...",                       │
│        "organization_id": "..." }  // NEW FIELD                 │
│                                                                 │
│    If Google OAuth:                                             │
│      POST /v1/auth/exchange_token                               │
│      { "code": "...", "redirect_url": "...", "type": "GOOGLE",  │
│        "organization_id": "..." }  // NEW FIELD                 │
│                                                                 │
│    Returns: org-scoped access + refresh tokens                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 8. Redirect to Console                                          │
│    /console (with organization context in token)                │
└─────────────────────────────────────────────────────────────────┘
```

## Architecture Changes

### 1. Database Schema

#### New Table: `email_verification_token`

```sql
CREATE TABLE email_verification_token
(
    email       VARCHAR(255) NOT NULL,
    token       VARCHAR(255) NOT NULL,
    created_at  BIGINT       NOT NULL,
    expires_at  BIGINT       NOT NULL,
    verified_at BIGINT DEFAULT NULL,
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    PRIMARY KEY (email),
    UNIQUE INDEX idx_token (token),
    INDEX       idx_expires_at (expires_at)
);
```

**Why separate table?**

- Keeps concerns separated: email verification vs password management
- Allows verification even for users without credentials yet
- Prevents confusion with password_reset_token field
- Can be cleaned up independently

#### Existing Tables (No Changes Needed)

**`account_v2`** - Already has:

- Composite PK: (email, organization_id)
- Organization role, environment roles
- Last seen tracking

**`account_credentials`** - Already has:

- email (PK)
- password_hash
- password_reset_token (for password reset/setup)
- password_reset_token_expires_at
- created_at, updated_at

**`organization`** - Already has:

- authentication_settings JSON column
- Structure: `{ enabled_types: [1, 2] }` where 1=GOOGLE, 2=PASSWORD

### 2. Protocol Buffer Definitions

#### New Messages in `proto/auth/service.proto`

```protobuf
// Email Verification Messages
message RequestMagicLinkRequest {
  string email = 1;
}

message RequestMagicLinkResponse {
  string message = 1; // Always: "If this email is registered, we've sent instructions."
}

message VerifyMagicLinkRequest {
  string token = 1;
}

message VerifyMagicLinkResponse {
  string email = 1; // Verified email
  repeated environment.Organization organizations = 2; // User's organizations
}
```

#### Modified Existing Messages

```protobuf
// Add organization_id field to existing SignInRequest
message SignInRequest {
  string email = 1;
  string password = 2;
  string organization_id = 3; // NEW - optional for backward compatibility
}

// Add organization_id field to existing ExchangeTokenRequest
message ExchangeTokenRequest {
  string code = 1;
  string redirect_url = 2;
  AuthType type = 3;
  string organization_id = 4; // NEW - optional for backward compatibility
}
```

#### Service Definition Updates

```protobuf
service AuthService {
  // ... existing methods ...

  // New Magic Link methods
  rpc RequestMagicLink(RequestMagicLinkRequest) returns (RequestMagicLinkResponse) {
    option (google.api.http) = {
      post: "/v1/auth/request_magic_link"
      body: "*"
    };
  }

  rpc VerifyMagicLink(VerifyMagicLinkRequest) returns (VerifyMagicLinkResponse) {
    option (google.api.http) = {
      post: "/v1/auth/verify_magic_link"
      body: "*"
    };
  }

  // Existing methods (SignIn and ExchangeToken) will be modified to accept
  // optional organization_id parameter for new flow, while maintaining
  // backward compatibility with old flow (auto-select first org if not provided)
}
```

### 3. Service Layer Implementation

#### New Storage Interface: `EmailVerificationStorage`

Location: `pkg/auth/storage/email_verification.go`

```go
type EmailVerificationStorage interface {
CreateVerificationToken(ctx context.Context, email, token string, expiresAt int64, ipAddress, userAgent string) error
GetVerificationToken(ctx context.Context, token string) (*EmailVerificationToken, error)
MarkVerified(ctx context.Context, token string) error
DeleteExpiredTokens(ctx context.Context, before int64) error
}

type EmailVerificationToken struct {
Email      string
Token      string
CreatedAt  int64
ExpiresAt  int64
VerifiedAt *int64
IPAddress  string
UserAgent  string
}
```

#### Updated: `authService` struct

Location: `pkg/auth/api/api.go`

```go
type authService struct {
// ... existing fields ...
emailVerificationStorage storage.EmailVerificationStorage // NEW
}
```

#### New Methods in `authService`

```go
// RequestMagicLink handles email verification initiation
func (s *authService) RequestMagicLink(
ctx context.Context,
req *authproto.RequestMagicLinkRequest,
) (*authproto.RequestMagicLinkResponse, error) {
// 1. Validate email format
// 2. Generate secure token (32 bytes)
// 3. Always create verification token record
// 4. Check if email exists (constant-time)
// 5. If exists: Send magic link email
// 6. Return same generic message regardless
}

// VerifyMagicLink handles magic link verification and returns org list
func (s *authService) VerifyMagicLink(
ctx context.Context,
req *authproto.VerifyMagicLinkRequest,
) (*authproto.VerifyMagicLinkResponse, error) {
// 1. Validate token exists and not expired
// 2. Get email from token
// 3. Mark token as verified
// 4. Call existing s.getOrganizationsByEmail(ctx, email, localizer)
// 5. Return email + organizations list
}
```

#### Modified Existing Methods

```go
// SignIn - Add organization_id parameter (optional for backward compat)
func (s *authService) SignIn(
ctx context.Context,
request *authproto.SignInRequest,
) (*authproto.SignInResponse, error) {
// Existing logic...

// NEW: If organization_id is provided:
if request.OrganizationId != "" {
// 1. Verify password using existing logic
// 2. Get organizations for email
// 3. Verify requested org is in the list
// 4. Get account for specific org
// 5. Verify PASSWORD auth is enabled for this org
// 6. Generate token for the specified org
} else {
// OLD: Existing behavior for backward compatibility
// Use first organization (current behavior)
}
}

// ExchangeToken - Add organization_id parameter (optional for backward compat)
func (s *authService) ExchangeToken(
ctx context.Context,
req *authproto.ExchangeTokenRequest,
) (*authproto.ExchangeTokenResponse, error) {
// Existing logic...

// NEW: If organization_id is provided:
if req.OrganizationId != "" {
// 1. Exchange OAuth code using existing googleAuthenticator
// 2. Get organizations for OAuth email
// 3. Verify requested org is in the list
// 4. Get account for specific org
// 5. Verify GOOGLE auth is enabled for this org
// 6. Generate token for the specified org
} else {
// OLD: Existing behavior for backward compatibility
// Use first organization (current behavior)
}
}
```

### 4. Email Service Integration

#### New Email Template

Location: `pkg/email/templates.go` (or similar)

```go
func (s *EmailService) SendMagicLinkEmail(
ctx context.Context,
email string,
magicLink string,
expiresIn time.Duration,
) error
```

**Email Template:**

```
Subject: Sign in to Bucketeer

Hello,

Click the link below to sign in to Bucketeer:

{{ magicLink }}

This link will expire in {{ expiresIn }} minutes.

If you didn't request this email, you can safely ignore it.

Thanks,
The Bucketeer Team
```

### 5. Frontend Requirements

The frontend implementation will need to support the following user flow:

#### User Journey

1. **Email Entry** (`/login`)
    - User enters email address
    - Shows generic "check your email" message after submission

2. **Email Verification** (`/auth/verify?token=...`)
    - User clicks magic link from email
    - Frontend calls `VerifyMagicLink` API
    - Receives verified email + list of organizations

3. **Organization Selection** (`/select-organization`)
    - Frontend displays organizations from previous step
    - User selects organization

4. **Authentication Method Selection** (`/auth/method-selection`)
    - Frontend displays available auth methods based on organization's `authentication_settings.enabled_types`
    - User chooses between Password or Google OAuth

5. **Final Authentication**
    - **Password**: Collect password, call `SignIn` with email, password, and organization_id
    - **Google OAuth**: Redirect to Google, callback calls `ExchangeToken` with code and organization_id

6. **Console Access** (`/console`)
    - Store access and refresh tokens
    - Redirect to main application

#### State Management

- Use `sessionStorage` to store email and organizations list after magic link verification (temporary state)
- Use `localStorage` for access/refresh tokens (persistent auth state)
- Clear session storage after successful authentication

### 6. Security Considerations

#### Rate Limiting

Implement at API gateway or service level:

- **Magic Link Requests**: 5 per IP per 15 minutes
- **Magic Link Verification**: 10 attempts per token (then invalidate)
- **Password Attempts**: 5 failed attempts per email per 15 minutes

#### Token Security

**Magic Link Token:**

- Cryptographically random (32+ bytes from crypto/rand)
- Single-use (mark as verified after first use)
- 15-minute expiry
- HTTPS only in production
- Stored in `email_verification_token` table

**Access/Refresh Tokens:**

- Existing implementation (already secure)
- Access: 24 hours
- Refresh: 7 days

#### Timing Attack Prevention

```go
func (s *authService) RequestMagicLink(
ctx context.Context,
req *authproto.RequestMagicLinkRequest,
) (*authproto.RequestMagicLinkResponse, error) {
email := req.Email

// Always create token record (constant operation)
token := auth.GenerateSecureToken()
expiresAt := time.Now().Add(15 * time.Minute).Unix()

_ = s.emailVerificationStorage.CreateVerificationToken(
ctx, email, token, expiresAt, getIP(ctx), getUserAgent(ctx),
)

// Check if user exists (database query - constant time for indexed email)
orgs, err := s.accountClient.GetMyOrganizationsByEmail(ctx,
&acproto.GetMyOrganizationsByEmailRequest{Email: email})

// Send email only if user exists, but don't reveal this
if err == nil && len(orgs.Organizations) > 0 {
magicLink := fmt.Sprintf("%s/auth/verify?token=%s",
s.config.BaseURL, token)
_ = s.emailService.SendMagicLinkEmail(ctx, email, magicLink, 15*time.Minute)
}

// Always return same response (constant time)
return &authproto.RequestMagicLinkResponse{
Message: "If this email is registered, we've sent instructions to your inbox.",
}, nil
}
```

#### CSRF Protection

- State parameter in OAuth flow (already implemented)
- SameSite cookies for pre-auth token (if using cookies instead of JWT in response)
- Validate origin/referer headers on state-changing operations

#### Session Management

- Email verification tokens: 15-minute expiry, single-use (mark as verified)
- Frontend session: Store email + organizations in sessionStorage after verification
- Cleanup job: Hourly deletion of expired/verified email_verification_token records

### 7. Edge Cases & Error Handling

#### User Has No Organizations

**Scenario**: Email verified but user has no org memberships

**Handling**:

```json
GET /v1/auth/organizations
Response 200: {
"organizations": []
}

Frontend shows: "You don't have access to any organizations. Contact your administrator."
```

#### Magic Link Expired

**Scenario**: User clicks magic link after 15 minutes

**Handling**: Return error, redirect to `/login?error=expired` with message:
"This link has expired. Please request a new one."

#### Magic Link Already Used

**Scenario**: User clicks magic link twice

**Handling**:

- Token is marked as verified after first use
- Second click: Check if token is already verified
- If verified recently (< 5 minutes): Return same email + org list (allow user to go back)
- If verified longer ago or expired: Show "Link already used, please request a new one"

#### Email Service Down

**Scenario**: Email service unavailable

**Handling**:

- Log error internally with high severity
- Still return generic success message to user (don't reveal service status)
- Alert ops team via monitoring
- User sees no email → can request again (rate limited)

#### User Changes Email Mid-Flow

**Scenario**: User requests link for email A, then requests for email B

**Handling**: Each request creates independent token. Both work until expiry. No conflict.

#### Organization Deleted Mid-Flow

**Scenario**: Org deleted between email verification and org selection

**Handling**: Filter out deleted/archived orgs in `GetOrganizationsForPreAuth`. If no orgs remain, show "no
organizations" error.

#### Auth Method Disabled Mid-Flow

**Scenario**: Admin disables password auth while user is in password entry screen

**Handling**:

- Validate enabled methods in `SignInWithOrg`
- Return error: `PERMISSION_DENIED - "Password authentication is not enabled for this organization"`
- Frontend redirects back to method selection or shows error

#### User Belongs to 50+ Organizations

**Scenario**: User has access to many orgs

**Handling**:

- Add pagination or search to org selection UI
- Backend: Return all orgs (no pagination needed for org list - typically < 10)
- Frontend: If > 10 orgs, add search box

#### Organization List Exposed in Frontend

**Scenario**: User's organization list is stored in sessionStorage

**Risk**: Minimal - organization names/IDs are not sensitive information

**Mitigation**:

- Still requires valid credentials (password or OAuth) to complete authentication
- HTTPS enforced (data encrypted in transit)
- Session storage cleared after successful login
- No sensitive data (roles, permissions) exposed until after authentication

### 8. Backward Compatibility

#### Existing APIs

Keep existing `SignIn` and `ExchangeToken` APIs for backward compatibility:

- Mark as deprecated in proto comments
- Add deprecation header to responses: `X-Deprecated: true`
- Log usage for monitoring migration progress
- Set deprecation timeline (e.g., 6 months)

#### Gradual Migration

- Old flow: Immediate token issuance for first org (existing behavior)
- New flow: Email verification → org selection → auth method → token
- Both flows coexist during transition period
- Feature flag controls which flow is shown to users

## References

- Bucketeer existing auth implementation: `pkg/auth/api/api.go`
- Bucketeer proto definitions: `proto/auth/service.proto`, `proto/environment/organization.proto`
