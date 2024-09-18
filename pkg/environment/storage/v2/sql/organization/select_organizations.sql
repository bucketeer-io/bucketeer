SELECT
    o.id,
    o.name,
    o.url_code,
    o.description,
    o.disabled,
    o.archived,
    o.trial,
    o.system_admin,
    o.created_at,
    o.updated_at,
    COUNT(DISTINCT p.id) AS projects,
    COUNT(DISTINCT e.id) AS environments,
    COUNT(DISTINCT a.email) AS users
FROM
    organization o
        LEFT JOIN project p ON o.id = p.organization_id
        LEFT JOIN environment_v2 e ON o.id = e.organization_id
        LEFT JOIN account_v2 a ON o.id = a.organization_id
GROUP BY
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
%s %s %s