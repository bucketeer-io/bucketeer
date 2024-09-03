UPDATE subscription SET
    updated_at = ?,
    disabled = ?,
    source_types = ?,
    recipient = ?,
    name = ?
WHERE
    id = ? AND
    environment_id = ?
