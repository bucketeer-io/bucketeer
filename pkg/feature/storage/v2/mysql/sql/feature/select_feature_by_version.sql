SELECT
    entity_data,
    COALESCE(feature.last_used_at, 0)
FROM audit_log
LEFT JOIN feature ON
    feature.id = audit_log.entity_id AND
    feature.environment_id = audit_log.environment_id
WHERE
    audit_log.environment_id = ? AND
    entity_type = 0 AND -- Feature entity type
    TRIM(entity_data) != '' AND
    entity_id = ? AND
    JSON_EXTRACT(entity_data, '$.version') = ?
LIMIT 1
