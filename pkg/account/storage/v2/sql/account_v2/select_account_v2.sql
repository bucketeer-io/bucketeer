SELECT
    email,
    name,
    avatar_image_url,
    organization_id,
    organization_role,
    environment_roles,
    disabled,
    created_at,
    updated_at,
    search_filters
FROM
    account_v2
WHERE
    email = ?
    AND organization_id = ?
