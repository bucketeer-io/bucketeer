SELECT
    a.email,
    a.name,
    a.first_name,
    a.last_name,
    a.language,
    a.avatar_image_url,
    a.avatar_file_type,
    a.avatar_image,
    a.tags,
    a.organization_id,
    a.organization_role,
    a.environment_roles,
    a.disabled,
    a.created_at,
    a.updated_at,
    a.last_seen,
    a.search_filters
FROM account_v2 AS a
INNER JOIN environment_v2 AS e
ON a.organization_id = e.organization_id
WHERE
    a.email IN (?) AND
    e.id = ?
