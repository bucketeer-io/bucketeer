-- Add new schema named "bucketeer"
CREATE DATABASE `bucketeer` CHARSET utf8 COLLATE utf8_general_ci;
-- Create "feature_last_used_info" table
CREATE TABLE `bucketeer`.`feature_last_used_info` (`id` varchar(255) NOT NULL, `feature_id` varchar(255) NOT NULL, `version` bigint NOT NULL, `last_used_at` bigint NOT NULL, `client_oldest_version` varchar(255) NOT NULL, `client_latest_version` varchar(255) NOT NULL, `created_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "push" table
CREATE TABLE `bucketeer`.`push` (`id` varchar(255) NOT NULL, `fcm_api_key` varchar(511) NOT NULL, `tags` json NOT NULL, `deleted` bool NOT NULL DEFAULT 0, `name` varchar(255) NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "account" table
CREATE TABLE `bucketeer`.`account` (`id` varchar(255) NOT NULL, `email` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `role` int NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `deleted` bool NOT NULL DEFAULT 0, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), UNIQUE INDEX `unique_email_environment_namespace` (`email`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "admin_audit_log" table
CREATE TABLE `bucketeer`.`admin_audit_log` (`id` varchar(255) NOT NULL, `timestamp` bigint NOT NULL, `entity_type` int NOT NULL, `entity_id` varchar(255) NOT NULL, `type` int NOT NULL, `event` json NOT NULL, `editor` json NOT NULL, `options` json NOT NULL, PRIMARY KEY (`id`), INDEX `idx_timestamp_desc` (`timestamp`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "webhook" table
CREATE TABLE `bucketeer`.`webhook` (`id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `description` text NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`), INDEX `idx_environment_namespace` (`environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "api_key" table
CREATE TABLE `bucketeer`.`api_key` (`id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `role` int NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "atlas_schema_revisions" table
CREATE TABLE `bucketeer`.`atlas_schema_revisions` (`version` varchar(255) NOT NULL, `description` varchar(255) NOT NULL, `type` bigint unsigned NOT NULL DEFAULT 2, `applied` bigint NOT NULL DEFAULT 0, `total` bigint NOT NULL DEFAULT 0, `executed_at` timestamp NULL, `execution_time` bigint NOT NULL, `error` longtext NULL, `error_stmt` longtext NULL, `hash` varchar(255) NOT NULL, `partial_hashes` json NULL, `operator_version` varchar(255) NOT NULL, PRIMARY KEY (`version`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "audit_log" table
CREATE TABLE `bucketeer`.`audit_log` (`id` varchar(255) NOT NULL, `timestamp` bigint NOT NULL, `entity_type` int NOT NULL, `entity_id` varchar(255) NOT NULL, `type` int NOT NULL, `event` json NOT NULL, `editor` json NOT NULL, `options` json NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `idx_entity_type_entity_id` (`entity_type`, `entity_id`), INDEX `idx_timestamp_desc` (`timestamp`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "tag" table
CREATE TABLE `bucketeer`.`tag` (`id` varchar(255) NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "subscription" table
CREATE TABLE `bucketeer`.`subscription` (`id` varchar(255) NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `source_types` json NOT NULL, `recipient` json NOT NULL, `name` varchar(255) NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "schema_migrations" table
CREATE TABLE `bucketeer`.`schema_migrations` (`version` bigint NOT NULL, `dirty` bool NOT NULL, PRIMARY KEY (`version`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "goal" table
CREATE TABLE `bucketeer`.`goal` (`id` varchar(255) NOT NULL, `name` varchar(511) NOT NULL, `description` text NOT NULL, `archived` bool NOT NULL DEFAULT 0, `deleted` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "experiment_result" table
CREATE TABLE `bucketeer`.`experiment_result` (`id` varchar(255) NOT NULL, `experiment_id` varchar(255) NOT NULL, `updated_at` bigint NOT NULL, `data` json NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "mau" table
CREATE TABLE `bucketeer`.`mau` (`user_id` varchar(255) NOT NULL, `yearmonth` varchar(6) NOT NULL, `source_id` varchar(30) NOT NULL, `event_count` int unsigned NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`user_id`, `yearmonth`, `source_id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "admin_account" table
CREATE TABLE `bucketeer`.`admin_account` (`id` varchar(255) NOT NULL, `email` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `role` int NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `deleted` bool NOT NULL DEFAULT 0, PRIMARY KEY (`id`), UNIQUE INDEX `unique_email` (`email`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "admin_subscription" table
CREATE TABLE `bucketeer`.`admin_subscription` (`id` varchar(255) NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `source_types` json NOT NULL, `recipient` json NOT NULL, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "organization" table
CREATE TABLE `bucketeer`.`organization` (`id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `url_code` varchar(255) NOT NULL, `description` text NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `archived` bool NOT NULL DEFAULT 0, `trial` bool NOT NULL DEFAULT 0, `system_admin` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unique_url_code` (`url_code`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "account_v2" table
CREATE TABLE `bucketeer`.`account_v2` (`email` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `avatar_image_url` varchar(255) NOT NULL, `organization_id` varchar(255) NOT NULL, `organization_role` int NOT NULL, `environment_roles` json NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, PRIMARY KEY (`email`, `organization_id`), INDEX `account_v2_foreign_organization_id` (`organization_id`), CONSTRAINT `account_v2_foreign_organization_id` FOREIGN KEY (`organization_id`) REFERENCES `bucketeer`.`organization` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "feature" table
CREATE TABLE `bucketeer`.`feature` (`id` varchar(255) NOT NULL, `name` varchar(511) NOT NULL, `description` text NOT NULL, `enabled` bool NOT NULL DEFAULT 0, `archived` bool NOT NULL DEFAULT 0, `deleted` bool NOT NULL DEFAULT 0, `evaluation_undelayable` bool NOT NULL DEFAULT 0, `ttl` int NOT NULL, `version` int NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `variations` json NOT NULL, `targets` json NOT NULL, `rules` json NOT NULL, `default_strategy` json NOT NULL, `off_variation` varchar(255) NOT NULL, `tags` json NOT NULL, `maintainer` varchar(255) NOT NULL, `environment_namespace` varchar(255) NOT NULL, `variation_type` int NOT NULL, `sampling_seed` varchar(255) NOT NULL DEFAULT "", `prerequisites` json NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "auto_ops_rule" table
CREATE TABLE `bucketeer`.`auto_ops_rule` (`id` varchar(255) NOT NULL, `feature_id` varchar(255) NOT NULL, `ops_type` int NOT NULL, `clauses` json NOT NULL, `triggered_at` bigint NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `deleted` bool NOT NULL DEFAULT 0, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `foreign_auto_ops_rule_feature_id_environment_namespace` (`feature_id`, `environment_namespace`), CONSTRAINT `foreign_auto_ops_rule_feature_id_environment_namespace` FOREIGN KEY (`feature_id`, `environment_namespace`) REFERENCES `bucketeer`.`feature` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "project" table
CREATE TABLE `bucketeer`.`project` (`id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `url_code` varchar(255) NOT NULL, `description` text NOT NULL, `disabled` bool NOT NULL DEFAULT 0, `trial` bool NOT NULL DEFAULT 0, `organization_id` varchar(255) NOT NULL, `creator_email` varchar(255) NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unique_url_code` (`url_code`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "environment" table
CREATE TABLE `bucketeer`.`environment` (`id` varchar(255) NOT NULL, `namespace` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `description` text NOT NULL, `deleted` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `project_id` varchar(255) NOT NULL, PRIMARY KEY (`id`), INDEX `foreign_project_id` (`project_id`), UNIQUE INDEX `unique_namespace` (`namespace`), CONSTRAINT `foreign_project_id` FOREIGN KEY (`project_id`) REFERENCES `bucketeer`.`project` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "environment_v2" table
CREATE TABLE `bucketeer`.`environment_v2` (`id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `url_code` varchar(255) NOT NULL, `description` text NOT NULL, `project_id` varchar(255) NOT NULL, `organization_id` varchar(255) NOT NULL, `archived` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unique_project_id_url_code` (`project_id`, `url_code`), CONSTRAINT `environment_v2_foreign_project_id` FOREIGN KEY (`project_id`) REFERENCES `bucketeer`.`project` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "experiment" table
CREATE TABLE `bucketeer`.`experiment` (`id` varchar(255) NOT NULL, `goal_id` varchar(255) NOT NULL, `feature_id` varchar(255) NOT NULL, `feature_version` int NOT NULL, `variations` json NOT NULL, `start_at` bigint NOT NULL, `stop_at` bigint NOT NULL, `stopped` bool NOT NULL DEFAULT 0, `stopped_at` bigint NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `archived` bool NOT NULL DEFAULT 0, `deleted` bool NOT NULL DEFAULT 0, `goal_ids` json NOT NULL, `name` varchar(511) NOT NULL, `description` text NOT NULL, `base_variation_id` varchar(255) NOT NULL, `status` int NOT NULL, `maintainer` varchar(255) NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `foreign_experiment_feature_id_environment_namespace` (`feature_id`, `environment_namespace`), CONSTRAINT `foreign_experiment_feature_id_environment_namespace` FOREIGN KEY (`feature_id`, `environment_namespace`) REFERENCES `bucketeer`.`feature` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "flag_trigger" table
CREATE TABLE `bucketeer`.`flag_trigger` (`id` varchar(255) NOT NULL, `feature_id` varchar(255) NOT NULL, `environment_namespace` varchar(255) NOT NULL, `type` int NOT NULL, `action` bool NOT NULL, `description` text NOT NULL, `trigger_count` int NOT NULL, `last_triggered_at` bigint NOT NULL, `uuid` varchar(512) NOT NULL DEFAULT "", `token` varchar(512) NOT NULL DEFAULT "", `disabled` bool NOT NULL DEFAULT 0, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `foreign_flag_triggers_feature_id_environment_namespace` (`feature_id`, `environment_namespace`), CONSTRAINT `flag_trigger_ibfk_1` FOREIGN KEY (`feature_id`, `environment_namespace`) REFERENCES `bucketeer`.`feature` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "ops_count" table
CREATE TABLE `bucketeer`.`ops_count` (`id` varchar(255) NOT NULL, `auto_ops_rule_id` varchar(255) NOT NULL, `clause_id` varchar(255) NOT NULL, `updated_at` bigint NOT NULL, `ops_event_count` bigint NOT NULL, `evaluation_count` bigint NOT NULL, `feature_id` varchar(255) NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `foreign_ops_count_auto_ops_rule_id_environment_namespace` (`auto_ops_rule_id`, `environment_namespace`), INDEX `foreign_ops_count_feature_id_environment_namespace` (`feature_id`, `environment_namespace`), CONSTRAINT `foreign_ops_count_auto_ops_rule_id_environment_namespace` FOREIGN KEY (`auto_ops_rule_id`, `environment_namespace`) REFERENCES `bucketeer`.`auto_ops_rule` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT, CONSTRAINT `foreign_ops_count_feature_id_environment_namespace` FOREIGN KEY (`feature_id`, `environment_namespace`) REFERENCES `bucketeer`.`feature` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "ops_progressive_rollout" table
CREATE TABLE `bucketeer`.`ops_progressive_rollout` (`id` varchar(255) NOT NULL, `feature_id` varchar(255) NOT NULL, `clause` json NOT NULL, `status` int NOT NULL, `stopped_by` int NOT NULL, `type` int NOT NULL, `stopped_at` bigint NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `foreign_progressive_rollout_feature_id_environment_namespace` (`feature_id`, `environment_namespace`), CONSTRAINT `foreign_progressive_rollout_feature_id_environment_namespace` FOREIGN KEY (`feature_id`, `environment_namespace`) REFERENCES `bucketeer`.`feature` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "segment" table
CREATE TABLE `bucketeer`.`segment` (`id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `description` text NOT NULL, `rules` json NOT NULL, `created_at` bigint NOT NULL, `updated_at` bigint NOT NULL, `version` bigint NOT NULL, `deleted` bool NOT NULL DEFAULT 0, `included_user_count` bigint NOT NULL, `excluded_user_count` bigint NOT NULL, `status` int NOT NULL, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "segment_user" table
CREATE TABLE `bucketeer`.`segment_user` (`id` varchar(511) NOT NULL, `segment_id` varchar(255) NOT NULL, `user_id` varchar(255) NOT NULL, `state` int NOT NULL, `deleted` bool NOT NULL DEFAULT 0, `environment_namespace` varchar(255) NOT NULL, PRIMARY KEY (`id`, `environment_namespace`), INDEX `foreign_segment_user_segment_id_environment_namespace` (`segment_id`, `environment_namespace`), CONSTRAINT `foreign_segment_user_segment_id_environment_namespace` FOREIGN KEY (`segment_id`, `environment_namespace`) REFERENCES `bucketeer`.`segment` (`id`, `environment_namespace`) ON UPDATE RESTRICT ON DELETE RESTRICT) CHARSET utf8mb4 COLLATE utf8mb4_bin;
