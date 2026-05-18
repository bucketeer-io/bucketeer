SELECT
    entity_data
FROM audit_log
WHERE
    environment_id = $1 AND
    entity_type = 0 AND
    TRIM(entity_data) != '' AND
    entity_id = $2 AND
    (entity_data::jsonb->>'version')::int = $3
LIMIT 1
