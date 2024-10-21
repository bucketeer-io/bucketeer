-- Modify "account_v2" table
ALTER TABLE `account_v2` ADD COLUMN `first_name` varchar(255) NULL, ADD COLUMN `last_name` varchar(255) NULL, ADD COLUMN `language` varchar(10) NULL, ADD COLUMN `last_seen` bigint NOT NULL DEFAULT 0;
