SELECT
    previous_entity_data
FROM audit_log
WHERE
    environment_id = ? AND
    entity_type = 0 AND -- Feature entity type
    TRIM(previous_entity_data) != '' AND
    entity_id = ? AND
    JSON_EXTRACT(previous_entity_data, '$.version') = ?
LIMIT 1
