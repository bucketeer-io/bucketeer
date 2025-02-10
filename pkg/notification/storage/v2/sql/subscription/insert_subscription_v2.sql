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
) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?)
