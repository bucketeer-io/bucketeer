-- Modify "account" table
ALTER TABLE `account` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "api_key" table
ALTER TABLE `api_key` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "audit_log" table
ALTER TABLE `audit_log` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "auto_ops_rule" table
ALTER TABLE `auto_ops_rule` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "experiment" table
ALTER TABLE `experiment` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "experiment_result" table
ALTER TABLE `experiment_result` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "feature" table
ALTER TABLE `feature` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "feature_last_used_info" table
ALTER TABLE `feature_last_used_info` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "flag_trigger" table
ALTER TABLE `flag_trigger` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "goal" table
ALTER TABLE `goal` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "mau" table
ALTER TABLE `mau` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "ops_count" table
ALTER TABLE `ops_count` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "ops_progressive_rollout" table
ALTER TABLE `ops_progressive_rollout` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "push" table
ALTER TABLE `push` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "segment" table
ALTER TABLE `segment` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "segment_user" table
ALTER TABLE `segment_user` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "subscription" table
ALTER TABLE `subscription` ADD COLUMN `environment_id` varchar(255) NULL;
-- Modify "tag" table
ALTER TABLE `tag` ADD COLUMN `environment_id` varchar(255) NULL;
