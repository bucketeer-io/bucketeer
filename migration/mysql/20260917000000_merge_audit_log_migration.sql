INSERT INTO audit_log (
  id, timestamp, entity_type, entity_id, type,
  event, editor, options, entity_data, previous_entity_data,
  environment_id, organization_id
)
SELECT id, timestamp, entity_type, entity_id, type,
       event, editor, options, entity_data, previous_entity_data,
       '', organization_id
FROM audit_log_migration
ON DUPLICATE KEY UPDATE audit_log.id = audit_log.id;
