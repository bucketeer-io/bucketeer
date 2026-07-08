UPDATE subscription SET
    updated_at = $1,
    disabled = $2,
    source_types = $3,
    recipient = $4,
    name = $5,
    feature_flag_tags = $6
WHERE
    id = $7 AND
    environment_id = $8
