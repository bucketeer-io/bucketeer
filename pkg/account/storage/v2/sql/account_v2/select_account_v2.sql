SELECT
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
    updated_at,
    last_seen,
    search_filters
FROM
    account_v2
WHERE
    email = ?
    AND organization_id = ?
