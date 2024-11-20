SELECT
    a.email,
    a.name,
    a.first_name,
    a.last_name,
    a.language,
    a.avatar_image_url,
    a.avatar_file_type,
    a.avatar_image,
    a.organization_id,
    a.organization_role,
    a.environment_roles,
    a.disabled,
    a.created_at,
    a.updated_at,
    a.last_seen,
    a.search_filters,
    o.id,
    o.name,
    o.url_code,
    o.description,
    o.disabled,
    o.archived,
    o.trial,
    o.system_admin,
    o.created_at,
    o.updated_at
FROM account_v2 AS a
INNER JOIN organization AS o
ON a.organization_id=o.id
WHERE email=?
