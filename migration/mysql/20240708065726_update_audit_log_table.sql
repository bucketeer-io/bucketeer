-- Modify "admin_audit_log" table
ALTER TABLE `admin_audit_log` ADD COLUMN `entity_data` text NOT NULL, ADD COLUMN `previous_entity_data` text NOT NULL;
-- Modify "audit_log" table
ALTER TABLE `audit_log` ADD COLUMN `entity_data` text NOT NULL, ADD COLUMN `previous_entity_data` text NOT NULL;
