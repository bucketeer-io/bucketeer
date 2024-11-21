-- Modify "api_key" table
ALTER TABLE `api_key` ADD COLUMN `api_key` varchar(255) NULL, ADD COLUMN `maintainer` varchar(255) NULL;
