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
    organization as org ON team.organization_id = org.id