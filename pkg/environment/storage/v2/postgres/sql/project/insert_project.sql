INSERT INTO project (
    id,
    name,
    url_code,
    description,
    disabled,
    trial,
    creator_email,
    organization_id,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
