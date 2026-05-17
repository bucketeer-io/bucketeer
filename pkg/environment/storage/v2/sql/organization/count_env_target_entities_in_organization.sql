SELECT
    COUNT(1)
FROM %s
LEFT JOIN environment_v2 ON %s.environment_id = environment_v2.id
WHERE environment_v2.organization_id = ?