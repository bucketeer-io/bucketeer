-- atlas:txmode none

-- Modify "audit_log" table: organization-level events are stored here with an
-- empty environment_id, so organization admins can read their own activity.
ALTER TABLE audit_log ADD COLUMN organization_id VARCHAR(255) NOT NULL DEFAULT '';
CREATE INDEX CONCURRENTLY idx_audit_log_organization_timestamp ON audit_log (organization_id, timestamp DESC);
