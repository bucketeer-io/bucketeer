INSERT INTO subscription (
    id,
    created_at,
    updated_at,
    disabled,
    source_types,
    recipient,
    name,
    environment_id
) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )
