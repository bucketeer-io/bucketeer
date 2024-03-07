UPDATE api_key SET
    name = ?,
    role = ?,
    disabled = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_namespace = ?
