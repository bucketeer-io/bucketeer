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
    environment_v2 AS e
ON
    a.organization_id = e.organization_id
WHERE
    a.email = ?
    AND e.id = ?