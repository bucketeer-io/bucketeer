SELECT email, 
    password_hash,
    created_at,
    updated_at
FROM account_credentials
WHERE email = ? 