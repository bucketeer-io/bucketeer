-- Modify "organization" table
-- First add the column as nullable
ALTER TABLE `organization` ADD COLUMN `authentication_settings` JSON NULL;

-- Set default authentication settings for all existing organizations (Google authentication only)
UPDATE `organization`
SET `authentication_settings` = '{"enabled_types": [1]}';

-- Make the column NOT NULL after setting values
ALTER TABLE `organization` MODIFY COLUMN `authentication_settings` JSON NOT NULL;