SELECT
    entity_data,
    COALESCE(feature.last_used_at, 0)
FROM audit_log
LEFT JOIN feature ON
    feature.id = audit_log.entity_id AND
    feature.environment_id = audit_log.environment_id
WHERE
    audit_log.environment_id = $1 AND
    entity_type = 0 AND
    TRIM(entity_data) != '' AND
    entity_id = $2 AND
    (entity_data::jsonb->>'version')::int = $3
LIMIT 1
