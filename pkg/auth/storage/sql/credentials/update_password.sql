UPDATE account_credentials
SET 
    password_hash = ?, 
    updated_at = ?
WHERE email = ? 