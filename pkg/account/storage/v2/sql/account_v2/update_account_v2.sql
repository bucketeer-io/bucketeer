UPDATE
    account_v2
SET
    name = ?,
    avatar_image_url = ?,
    organization_role = ?,
    environment_roles = ?,
    disabled = ?,
    updated_at = ?
WHERE
    email = ?
    AND organization_id = ?