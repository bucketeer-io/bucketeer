-- Modify "account_v2" table
ALTER TABLE `account_v2`
    MODIFY COLUMN `name` varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN `first_name` varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN `last_name` varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN `language` varchar(10) NOT NULL DEFAULT '',
    ADD COLUMN `last_seen` bigint NOT NULL DEFAULT 0,
    ADD COLUMN `avatar_file_type` varchar(50) NOT NULL DEFAULT '',
    ADD COLUMN `avatar_image` mediumblob NULL;