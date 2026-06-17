SELECT
    team.id,
    team.name,
    team.description,
    team.created_at,
    team.updated_at,
    team.organization_id,
    org.name as organization_name
FROM
    team
JOIN
    organization org ON team.organization_id = org.id
WHERE
    team.name = $1 AND
    team.organization_id = $2
