-- Modify "audit_log" table: organization-level events are stored here with an
-- empty environment_id, so organization admins can read their own activity.
ALTER TABLE `audit_log` ADD COLUMN `organization_id` varchar(255) NOT NULL DEFAULT '', ALGORITHM=INSTANT;
CREATE INDEX `idx_organization_id_timestamp_desc` ON `audit_log` (`organization_id`, `timestamp` DESC);
