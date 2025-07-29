UPDATE account_credentials
SET 
    password_reset_token = ?,
    password_reset_token_expires_at = ?,
    updated_at = ?
WHERE email = ? 