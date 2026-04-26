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
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    organization_role = VALUES(organization_role),
    environment_roles = VALUES(environment_roles),
    updated_at = VALUES(updated_at)
