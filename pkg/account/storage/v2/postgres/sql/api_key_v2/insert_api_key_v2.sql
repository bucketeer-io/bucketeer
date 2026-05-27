INSERT INTO api_key (
    id,
    name,
    role,
    disabled,
    created_at,
    updated_at,
    environment_id,
    api_key,
    maintainer,
    description
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
