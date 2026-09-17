-- Merge the resolved admin_audit_log history from the audit_log_migration working table
-- into audit_log (follows 20260907000000). Organization-level rows land with an empty
-- environment_id; rows still unresolved are system-level and keep an empty organization_id.
-- No-op on conflict and where the working table is empty, so the file is safe to re-apply.
-- admin_audit_log and audit_log_migration are kept and dropped in a later release.

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
