DELETE FROM email_verification_token
WHERE expires_at < ?
OR (verified_at IS NOT NULL AND verified_at < ?)
