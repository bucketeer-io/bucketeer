-- Populate environment_id in feature table
UPDATE `feature` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in account table
UPDATE `account` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in api_key table
UPDATE `api_key` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in audit_log table
UPDATE `audit_log` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in auto_ops_rule table
UPDATE `auto_ops_rule` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in experiment table
UPDATE `experiment` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in experiment_result table
UPDATE `experiment_result` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in feature_last_used_info table
UPDATE `feature_last_used_info` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in flag_trigger table
UPDATE `flag_trigger` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in goal table
UPDATE `goal` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in mau table
UPDATE `mau` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in ops_count table
UPDATE `ops_count` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in ops_progressive_rollout table
UPDATE `ops_progressive_rollout` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in push table
UPDATE `push` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in segment table
UPDATE `segment` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in segment_user table
UPDATE `segment_user` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in subscription table
UPDATE `subscription` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;
-- Populate environment_id in tag table
UPDATE `tag` SET `environment_id` = `environment_namespace` WHERE `environment_id` IS NULL;