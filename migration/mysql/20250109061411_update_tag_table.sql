ALTER TABLE `tag` ADD COLUMN `name` VARCHAR(255) NOT NULL AFTER `id`;

-- Populate name with values from id
UPDATE `tag` SET name = id;
