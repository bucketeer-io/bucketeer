-- Add the name column
ALTER TABLE `tag` ADD COLUMN `name` VARCHAR(255) NOT NULL AFTER `id`;

-- Populate name with values from id
UPDATE `tag` SET name = id;

-- Make the name unique
ALTER TABLE `tag` ADD UNIQUE (`name`, `environment_id`, `entity_type`);

-- Drop the `environment_id` from the primary key
ALTER TABLE `tag` DROP PRIMARY KEY;
ALTER TABLE `tag` ADD PRIMARY KEY (`id`);
