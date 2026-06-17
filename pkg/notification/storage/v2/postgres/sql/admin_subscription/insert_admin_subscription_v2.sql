INSERT INTO admin_subscription (
    id,
    created_at,
    updated_at,
    disabled,
    source_types,
    recipient,
    name
) VALUES ($1, $2, $3, $4, $5, $6, $7)
