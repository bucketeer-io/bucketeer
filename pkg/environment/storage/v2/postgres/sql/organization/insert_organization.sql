INSERT INTO organization (
    id,
    name,
    owner_email,
    url_code,
    description,
    disabled,
    archived,
    trial,
    system_admin,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
