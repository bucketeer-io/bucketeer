UPDATE api_key SET
    name = ?,
    role = ?,
    disabled = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_namespace = ?