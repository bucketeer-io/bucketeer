CREATE TABLE IF NOT EXISTS `account` (
    `id` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `role` INT(11) NOT NULL,
    `disabled` TINYINT(1) NOT NULL DEFAULT '0',
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    UNIQUE INDEX `unique_email_environment_namespace` (`email` ASC, `environment_namespace` ASC)
);

CREATE TABLE IF NOT EXISTS `admin_account` (
    `id` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `role` INT(11) NOT NULL,
    `disabled` TINYINT(1) NOT NULL DEFAULT '0',
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `unique_email` (`email`)
);

CREATE TABLE IF NOT EXISTS `admin_audit_log` (
    `id` VARCHAR(255) NOT NULL,
    `timestamp` BIGINT(20) NOT NULL,
    `entity_type` INT(11) NOT NULL,
    `entity_id`  VARCHAR(255) NOT NULL,
    `type` INT(11) NOT NULL,
    `event` JSON NOT NULL,
    `editor` JSON NOT NULL,
    `options` JSON NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_timestamp_desc` (`timestamp` DESC)
);

CREATE TABLE IF NOT EXISTS `admin_subscription` (
    `id` VARCHAR(255) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `disabled` TINYINT(1) NOT NULL DEFAULT '0',
    `source_types` JSON NOT NULL,
    `recipient` JSON NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `api_key` (
    `id` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `role` INT(11) NOT NULL,
    `disabled` TINYINT(1) NOT NULL DEFAULT '0',
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `audit_log` (
    `id` VARCHAR(255) NOT NULL,
    `timestamp` BIGINT(20) NOT NULL,
    `entity_type` INT(11) NOT NULL,
    `entity_id`  VARCHAR(255) NOT NULL,
    `type` INT(11) NOT NULL,
    `event` JSON NOT NULL,
    `editor` JSON NOT NULL,
    `options` JSON NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    INDEX `idx_entity_type_entity_id` (`entity_type` ASC, `entity_id` ASC),
    INDEX `idx_timestamp_desc` (`timestamp` DESC)
);

CREATE TABLE IF NOT EXISTS `auto_ops_rule` (
    `id` VARCHAR(255) NOT NULL,
    `feature_id` VARCHAR(255) NOT NULL,
    `ops_type` INT(11) NOT NULL,
    `clauses` JSON NOT NULL,
    `triggered_at` BIGINT(20) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    CONSTRAINT `foreign_auto_ops_rule_feature_id_environment_namespace`
    FOREIGN KEY (`feature_id`,`environment_namespace`)
    REFERENCES `feature` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS `environment` (
    `id` VARCHAR(255) NOT NULL,
    `namespace` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `project_id` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `unique_namespace` (`namespace` ASC),
    CONSTRAINT `foreign_project_id`
    FOREIGN KEY (`project_id`)
    REFERENCES `project` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS `experiment` (
    `id` VARCHAR(255) NOT NULL,
    `goal_id` VARCHAR(255) NOT NULL,
    `feature_id` VARCHAR(255) NOT NULL,
    `feature_version` INT(11) NOT NULL,
    `variations` JSON NOT NULL,
    `start_at` BIGINT(20) NOT NULL,
    `stop_at` BIGINT(20) NOT NULL,
    `stopped` TINYINT(1) NOT NULL DEFAULT '0',
    `stopped_at` BIGINT(20) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `archived` TINYINT(1) NOT NULL DEFAULT '0',
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `goal_ids` JSON NOT NULL,
    `name` VARCHAR(511) NOT NULL,
    `description` TEXT NOT NULL,
    `base_variation_id` VARCHAR(255) NOT NULL,
    `status` INT(11) NOT NULL,
    `maintainer` VARCHAR(255) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    CONSTRAINT `foreign_experiment_feature_id_environment_namespace`
    FOREIGN KEY (`feature_id`, `environment_namespace`)
    REFERENCES `feature` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS `experiment_result` (
    `id` VARCHAR(255) NOT NULL,
    `experiment_id` VARCHAR(255) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `data` JSON NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `feature` (
    `id` VARCHAR(255) NOT NULL,
    `name` VARCHAR(511) NOT NULL,
    `description` TEXT NOT NULL,
    `enabled` TINYINT(1) NOT NULL DEFAULT '0',
    `archived` TINYINT(1) NOT NULL DEFAULT '0',
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `evaluation_undelayable` TINYINT(1) NOT NULL DEFAULT '0',
    `ttl` INT(11) NOT NULL,
    `version` INT(11) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `variations` JSON NOT NULL,
    `targets` JSON NOT NULL,
    `rules` JSON NOT NULL,
    `default_strategy` JSON NOT NULL,
    `off_variation` VARCHAR(255) NOT NULL,
    `tags` JSON NOT NULL,
    `maintainer` VARCHAR(255) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    `variation_type` INT(11) NOT NULL,
    `sampling_seed` VARCHAR(255) NOT NULL DEFAULT '',
    `prerequisites` JSON DEFAULT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `feature_last_used_info` (
    `id` VARCHAR(255) NOT NULL,
    `feature_id` VARCHAR(255) NOT NULL,
    `version` BIGINT(11) NOT NULL,
    `last_used_at` BIGINT(20) NOT NULL,
    `client_oldest_version` VARCHAR(255) NOT NULL,
    `client_latest_version` VARCHAR(255) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `goal` (
    `id` VARCHAR(255) NOT NULL,
    `name` VARCHAR(511) NOT NULL,
    `description` TEXT NOT NULL,
    `archived` TINYINT(1) NOT NULL DEFAULT '0',
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `mau` (
    `user_id` VARCHAR(255) NOT NULL,
    `yearmonth` VARCHAR(6) NOT NULL,
    `source_id` VARCHAR(30) NOT NULL,
    `event_count` INT(11) UNSIGNED NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`environment_namespace`, `yearmonth`, `source_id`, `user_id`)
);

CREATE TABLE IF NOT EXISTS `ops_count` (
    `id` VARCHAR(255) NOT NULL,
    `auto_ops_rule_id` VARCHAR(255) NOT NULL,
    `clause_id` VARCHAR(255) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `ops_event_count` BIGINT(20) NOT NULL,
    `evaluation_count` BIGINT(20) NOT NULL,
    `feature_id` VARCHAR(255) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    CONSTRAINT `foreign_ops_count_auto_ops_rule_id_environment_namespace`
    FOREIGN KEY (`auto_ops_rule_id`, `environment_namespace`)
    REFERENCES `auto_ops_rule` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
    CONSTRAINT `foreign_ops_count_feature_id_environment_namespace`
    FOREIGN KEY (`feature_id`, `environment_namespace`)
    REFERENCES `feature` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS `ops_progressive_rollout` (
    `id` VARCHAR(255) NOT NULL,
    `feature_id` VARCHAR(255) NOT NULL,
    `clause` JSON NOT NULL,
    `status` INT(11) NOT NULL,
    `type` INT(11) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    CONSTRAINT `foreign_progressive_rollout_feature_id_environment_namespace`
    FOREIGN KEY (`feature_id`, `environment_namespace`)
    REFERENCES `feature` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS `project` (
    `id` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `disabled` TINYINT(1) NOT NULL DEFAULT '0',
    `trial` TINYINT(1) NOT NULL DEFAULT '0',
    `creator_email` VARCHAR(255) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `push` (
    `id` VARCHAR(255) NOT NULL,
    `fcm_api_key` VARCHAR(511) NOT NULL,
    `tags` JSON NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `name` VARCHAR(255) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `segment` (
    `id` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `rules` JSON NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `version` BIGINT(20) NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `included_user_count` BIGINT(20) NOT NULL,
    `excluded_user_count` BIGINT(20) NOT NULL,
    `status` INT(11) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);

CREATE TABLE IF NOT EXISTS `segment_user` (
    `id` VARCHAR(511) NOT NULL,
    `segment_id` VARCHAR(255) NOT NULL,
    `user_id` VARCHAR(255) NOT NULL,
    `state` INT(11) NOT NULL,
    `deleted` TINYINT(1) NOT NULL DEFAULT '0',
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`),
    CONSTRAINT `foreign_segment_user_segment_id_environment_namespace`
    FOREIGN KEY (`segment_id`, `environment_namespace`)
    REFERENCES `segment` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);
















CREATE TABLE IF NOT EXISTS `subscription` (
    `id` VARCHAR(255) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `disabled` TINYINT(1) NOT NULL DEFAULT '0',
    `source_types` JSON NOT NULL,
    `recipient` JSON NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);














CREATE TABLE IF NOT EXISTS `tag` (
    `id` VARCHAR(255) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`, `environment_namespace`)
);
















CREATE TABLE IF NOT EXISTS `webhook` (
    `id` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `environment_namespace` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_environment_namespace` (`environment_namespace`)
);
