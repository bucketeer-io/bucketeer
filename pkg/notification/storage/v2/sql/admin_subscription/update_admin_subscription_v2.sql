UPDATE admin_subscription SET
    updated_at = ?,
    disabled = ?,
    source_types = ?,
    recipient = ?,
    name = ?
WHERE
    id = ?
