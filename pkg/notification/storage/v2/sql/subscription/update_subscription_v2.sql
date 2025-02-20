UPDATE subscription SET
    updated_at = ?,
    disabled = ?,
    source_types = ?,
    recipient = ?,
    name = ?,
    feature_flag_tags = ?
WHERE
    id = ? AND
    environment_id = ?
