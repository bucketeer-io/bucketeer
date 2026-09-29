UPDATE account_credentials
SET 
    password_reset_token = NULL, 
    password_reset_token_expires_at = NULL, 
    updated_at = ?
WHERE password_reset_token = ? 