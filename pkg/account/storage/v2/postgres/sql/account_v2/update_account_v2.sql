UPDATE
    account_v2
SET
    name = $1,
    first_name = $2,
    last_name = $3,
    language = $4,
    avatar_image_url = $5,
    avatar_file_type = $6,
    avatar_image = $7,
    tags = $8,
    teams = $9,
    organization_role = $10,
    environment_roles = $11,
    disabled = $12,
    updated_at = $13,
    last_seen = $14,
    search_filters = $15
WHERE
    email = $16
    AND organization_id = $17
