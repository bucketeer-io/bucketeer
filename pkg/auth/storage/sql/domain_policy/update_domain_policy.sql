UPDATE domain_auth_policy
SET auth_policy = ?,
    enabled = ?,
    updated_at = ?
WHERE domain = ?