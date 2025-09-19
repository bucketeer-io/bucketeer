-- Modify "organization" table
ALTER TABLE `organization` ADD COLUMN `password_authentication_enabled` BOOLEAN NOT NULL DEFAULT TRUE;