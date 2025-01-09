-- Set the default entity type to 1 (FEATURE_FLAG). This will be updated to zero after the migration is finished.
ALTER TABLE `tag` ADD COLUMN `entity_type` INT NOT NULL DEFAULT 1 AFTER `environment_id`;
