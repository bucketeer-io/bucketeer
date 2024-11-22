-- Modify "api_key" table
ALTER TABLE `api_key` ADD COLUMN `api_key` varchar(255) NOT NULL DEFAULT "", ADD COLUMN `maintainer` varchar(255) NOT NULL DEFAULT "", ADD COLUMN `description` varchar(255) NOT NULL DEFAULT "";
