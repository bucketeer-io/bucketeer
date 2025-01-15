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
    organization_id,
    organization_role,
    environment_roles,
    disabled,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
