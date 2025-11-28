UPDATE email_verification_token
SET verified_at = ?
WHERE token = ?
AND verified_at IS NULL
