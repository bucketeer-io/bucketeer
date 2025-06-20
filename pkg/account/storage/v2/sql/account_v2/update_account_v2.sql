UPDATE
    account_v2
SET
    name = ?,
    first_name = ?,
    last_name = ?,
    language = ?,
    avatar_image_url = ?,
    avatar_file_type = ?,
    avatar_image = ?,
    tags = ?,
    teams = ?,
    organization_role = ?,
    environment_roles = ?,
    disabled = ?,
    updated_at = ?,
    last_seen = ?,
    search_filters = ?
WHERE
    email = ?
    AND organization_id = ?
