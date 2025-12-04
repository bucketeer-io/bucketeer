INSERT INTO email_verification_token (
    email,
    token,
    created_at,
    expires_at,
    verified_at,
    ip_address,
    user_agent
) VALUES (?, ?, ?, ?, NULL, ?, ?)
ON DUPLICATE KEY UPDATE
    token = VALUES(token),
    created_at = VALUES(created_at),
    expires_at = VALUES(expires_at),
    verified_at = NULL,
    ip_address = VALUES(ip_address),
    user_agent = VALUES(user_agent)
