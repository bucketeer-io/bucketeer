SELECT
    a.email,
    a.name,
    a.avatar_image_url,
    a.organization_id,
    a.organization_role,
    a.environment_roles,
    a.disabled,
    a.created_at,
    a.updated_at
FROM
    account_v2 AS a
INNER JOIN
    organization AS o
ON
    a.organization_id = o.id
WHERE
    o.system_admin = 1
    AND a.email = ?