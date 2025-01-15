ALTER TABLE `account_v2` ADD COLUMN `tags` json NOT NULL AFTER `avatar_image`;

-- Populate name with values from id
UPDATE `account_v2` SET tags = '[]';
