INSERT INTO account_v2 (
    email,
    name,
    avatar_image_url,
    tags,
    organization_id,
    organization_role,
    environment_roles,
    disabled,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS new
ON DUPLICATE KEY UPDATE
    name = new.name,
    avatar_image_url = new.avatar_image_url,
    tags = new.tags,
    organization_role = new.organization_role,
    environment_roles = new.environment_roles,
    disabled = new.disabled,
    updated_at = new.updated_at
