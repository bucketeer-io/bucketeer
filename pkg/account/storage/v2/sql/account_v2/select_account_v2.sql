SELECT
    email,
    name,
    first_name,
    last_name,
    language,
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
