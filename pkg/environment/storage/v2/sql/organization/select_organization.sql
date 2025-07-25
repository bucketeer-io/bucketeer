SELECT
    organization.id,
    organization.name,
    organization.owner_email,
    organization.url_code,
    organization.description,
    organization.disabled,
    organization.archived,
    organization.trial,
    organization.system_admin,
    organization.created_at,
    organization.updated_at,
    (SELECT COUNT(DISTINCT id) FROM project WHERE organization_id = organization.id) AS projects,
    (SELECT COUNT(DISTINCT id) FROM environment_v2 WHERE organization_id = organization.id) AS environments,
    (SELECT COUNT(DISTINCT email) FROM account_v2 WHERE organization_id = organization.id) AS users
FROM organization
WHERE organization.id = ?