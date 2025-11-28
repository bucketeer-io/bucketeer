SELECT
    email,
    token,
    created_at,
    expires_at,
    verified_at,
    ip_address,
    user_agent
FROM email_verification_token
WHERE token = ?
