SELECT 
    password_reset_token, 
    email, 
    password_reset_token_expires_at, 
    created_at
FROM account_credentials
WHERE password_reset_token = ? AND password_reset_token IS NOT NULL 