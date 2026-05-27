INSERT INTO account_v2 (
    email,
    name,
    first_name,
    last_name,
    language,
    avatar_image_url,
    avatar_file_type,
    avatar_image,
    tags,
    teams,
    organization_id,
    organization_role,
    environment_roles,
    disabled,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
