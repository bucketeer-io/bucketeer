SELECT
    organization.id,
    organization.name,
    organization.url_code,
    organization.description,
    organization.disabled,
    organization.archived,
    organization.trial,
    organization.system_admin,
    organization.created_at,
    organization.updated_at,
    COUNT(DISTINCT project.id) AS projects,
    COUNT(DISTINCT environment_v2.id) AS environments,
    COUNT(DISTINCT account_v2.email) AS users
FROM
    organization
LEFT JOIN project ON organization.id = project.organization_id
LEFT JOIN environment_v2 ON organization.id = environment_v2.organization_id
LEFT JOIN account_v2 ON organization.id = account_v2.organization_id
%s
GROUP BY
    organization.id,
    organization.name,
    organization.url_code,
    organization.description,
    organization.disabled,
    organization.archived,
    organization.trial,
    organization.system_admin,
    organization.created_at,
    organization.updated_at
%s
%s