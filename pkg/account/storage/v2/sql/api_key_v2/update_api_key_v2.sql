UPDATE api_key SET
    name = ?,
    role = ?,
    disabled = ?,
    maintainer = ?,
    description = ?,
    updated_at = ?,
    last_used_at = ?
WHERE
    id = ? AND
    environment_id = ?