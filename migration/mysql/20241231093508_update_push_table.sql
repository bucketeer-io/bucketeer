ALTER TABLE `push` ADD COLUMN `disabled` tinyint(1) NOT NULL DEFAULT '0' AFTER `tags`;
