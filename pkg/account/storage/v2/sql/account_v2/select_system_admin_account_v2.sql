SELECT
    a.email,
    a.name,
    a.first_name,
    a.last_name,
    a.language,
    a.avatar_image_url,
    a.organization_id,
    a.organization_role,
    a.environment_roles,
    a.disabled,
    a.created_at,
    a.updated_at,
    a.search_filters
FROM
    account_v2 AS a
INNER JOIN
    organization AS o
ON
    a.organization_id = o.id
WHERE
    o.system_admin = 1
    AND a.email = ?
