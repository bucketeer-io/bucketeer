INSERT INTO subscription (
    id,
    created_at,
    updated_at,
    disabled,
    source_types,
    recipient,
    name,
    feature_flag_tags,
    environment_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
