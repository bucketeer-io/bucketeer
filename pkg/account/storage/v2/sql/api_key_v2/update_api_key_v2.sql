UPDATE api_key SET
    name = ?,
    role = ?,
    disabled = ?,
    maintainer = ?,
    description = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_id = ?