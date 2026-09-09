-- Resolve the organization for the frozen admin_audit_log history in a working table.
-- The merge into audit_log is a follow-up migration, after the resolved_by counts are reviewed.
-- A row can still reach admin_audit_log after this snapshot (the persister keeps a temporary
-- fallback for events published by producers that predate organization_id); the follow-up
-- migration delta-copies such late rows and re-runs these rules, which are idempotent.
-- No-op where admin_audit_log is empty; each UPDATE only touches rows still unresolved (organization_id = '').

CREATE TABLE audit_log_migration LIKE admin_audit_log;

ALTER TABLE audit_log_migration
  ADD COLUMN organization_id varchar(255) NOT NULL DEFAULT '',
  ADD COLUMN resolved_by varchar(32) NOT NULL DEFAULT '';

INSERT INTO audit_log_migration
  (id, timestamp, entity_type, entity_id, type, event, editor, options,
   entity_data, previous_entity_data)
SELECT id, timestamp, entity_type, entity_id, type, event, editor, options,
       entity_data, previous_entity_data
FROM admin_audit_log;

-- 1. Organization in the entity snapshot (JSON_VALID required: entity_data is text, often empty).
UPDATE audit_log_migration
SET organization_id = JSON_UNQUOTE(JSON_EXTRACT(entity_data, '$.organization_id')),
    resolved_by     = '1-entity-data'
WHERE organization_id = ''
  AND JSON_VALID(entity_data)
  AND COALESCE(JSON_UNQUOTE(JSON_EXTRACT(entity_data, '$.organization_id')), '') <> '';

-- 2. Organization in the previous snapshot.
UPDATE audit_log_migration
SET organization_id = JSON_UNQUOTE(JSON_EXTRACT(previous_entity_data, '$.organization_id')),
    resolved_by     = '2-previous-entity-data'
WHERE organization_id = ''
  AND JSON_VALID(previous_entity_data)
  AND COALESCE(JSON_UNQUOTE(JSON_EXTRACT(previous_entity_data, '$.organization_id')), '') <> '';

-- 3. The row is itself an organization (entity_type 15), so entity_id is the organization ID.
UPDATE audit_log_migration
SET organization_id = COALESCE(
      NULLIF(CASE WHEN JSON_VALID(entity_data)
                  THEN JSON_UNQUOTE(JSON_EXTRACT(entity_data, '$.id')) END, ''),
      entity_id),
    resolved_by = '3-organization-row'
WHERE organization_id = '' AND entity_type = 15;

-- 4. Environment row (entity_type 6): entity_id is an environment ID.
UPDATE audit_log_migration m
  JOIN environment_v2 e ON e.id = m.entity_id
SET m.organization_id = e.organization_id,
    m.resolved_by     = '4-environment'
WHERE m.organization_id = '' AND m.entity_type = 6;

-- 5. Project row (entity_type 12): entity_id is a project ID.
UPDATE audit_log_migration m
  JOIN project p ON p.id = m.entity_id
SET m.organization_id = p.organization_id,
    m.resolved_by     = '5-project'
WHERE m.organization_id = '' AND m.entity_type = 12;

-- 6. Type 309 (environment roles changed): find environment IDs in the decoded payload.
-- The 0x0A+length prefix anchors the match to a protobuf string field;
-- accepted only when all matches agree on one organization.
UPDATE audit_log_migration m
  JOIN (
    SELECT m2.id, MIN(e.organization_id) AS organization_id
    FROM audit_log_migration m2
    JOIN environment_v2 e
      ON LOCATE(CONCAT(UNHEX('0A'), UNHEX(LPAD(HEX(LENGTH(e.id)), 2, '0')), e.id),
                FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(m2.event, '$.value')))) > 0
    WHERE m2.organization_id = '' AND m2.type = 309
    GROUP BY m2.id
    HAVING COUNT(DISTINCT e.organization_id) = 1
  ) r ON r.id = m.id
SET m.organization_id = r.organization_id,
    m.resolved_by     = '6-role-payload';

-- 7. Account rows (entity_type 3, and 9 — some account events were mislabeled PUSH):
-- entity_id is an email; only when it maps to exactly one organization. account_v2 is
-- current state and deletion removes rows, so an email that moved organizations can map
-- uniquely to the wrong one; decline when rows of this log already resolved the same
-- email to a different organization.
UPDATE audit_log_migration m
  JOIN (
    SELECT email, MIN(organization_id) AS organization_id
    FROM account_v2
    GROUP BY email
    HAVING COUNT(DISTINCT organization_id) = 1
  ) av ON av.email = m.entity_id
  LEFT JOIN (
    SELECT DISTINCT entity_id, organization_id
    FROM audit_log_migration
    WHERE organization_id <> ''
  ) conflict ON conflict.entity_id = m.entity_id AND conflict.organization_id <> av.organization_id
SET m.organization_id = av.organization_id,
    m.resolved_by     = '7-account'
WHERE m.organization_id = '' AND m.entity_type IN (3, 9)
  AND conflict.entity_id IS NULL;

-- 8. Same email already resolved on another row of this log (unique org only). Must run last.
UPDATE audit_log_migration m
  JOIN (
    SELECT entity_id, MIN(organization_id) AS organization_id
    FROM audit_log_migration
    WHERE organization_id <> '' AND entity_type IN (3, 9)
    GROUP BY entity_id
    HAVING COUNT(DISTINCT organization_id) = 1
  ) x ON x.entity_id = m.entity_id
SET m.organization_id = x.organization_id,
    m.resolved_by     = '8-email-elsewhere'
WHERE m.organization_id = '' AND m.entity_type IN (3, 9);
