UPDATE api_key SET
    name = ?,
    role = ?,
    disabled = ?,
    description = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_id = ?