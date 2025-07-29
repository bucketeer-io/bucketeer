SELECT email, 
    password_hash,
    password_reset_token,
    password_reset_token_expires_at,
    created_at,
    updated_at
FROM account_credentials
WHERE email = ? 