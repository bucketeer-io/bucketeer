ALTER TABLE `subscription` ADD COLUMN `feature_flag_tags` JSON NOT NULL AFTER `name`;

-- Populate feature_flag_tags with an empty array
UPDATE `subscription` SET `feature_flag_tags` = '[]';
