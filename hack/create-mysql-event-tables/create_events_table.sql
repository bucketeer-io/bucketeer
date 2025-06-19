-- Create Evaluation Event Table
CREATE TABLE IF NOT EXISTS `evaluation_event` (
    `id` VARCHAR(255) NOT NULL,
    `environment_id` VARCHAR(255) NOT NULL,
    `timestamp` DATETIME(6) NOT NULL,
    `feature_id` VARCHAR(255) NOT NULL,
    `feature_version` INT NOT NULL,
    `user_id` VARCHAR(255) NOT NULL,
    `user_data` JSON,
    `variation_id` VARCHAR(255) NOT NULL,
    `reason` TEXT,
    `tag` VARCHAR(255),
    `source_id` VARCHAR(255),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_evaluation_environment_id` (`environment_id`),
    INDEX `idx_evaluation_timestamp` (`timestamp`),
    INDEX `idx_evaluation_feature_id` (`feature_id`),
    INDEX `idx_evaluation_user_id` (`user_id`),
    INDEX `idx_evaluation_variation_id` (`variation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Goal Event Table
CREATE TABLE IF NOT EXISTS `goal_event` (
    `id` VARCHAR(255) NOT NULL,
    `environment_id` VARCHAR(255) NOT NULL,
    `timestamp` DATETIME(6) NOT NULL,
    `goal_id` VARCHAR(255) NOT NULL,
    `value` FLOAT,
    `user_id` VARCHAR(255) NOT NULL,
    `user_data` JSON,
    `tag` VARCHAR(255),
    `source_id` VARCHAR(255),
    `feature_id` VARCHAR(255),
    `feature_version` INT,
    `variation_id` VARCHAR(255),
    `reason` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_goal_environment_id` (`environment_id`),
    INDEX `idx_goal_timestamp` (`timestamp`),
    INDEX `idx_goal_goal_id` (`goal_id`),
    INDEX `idx_goal_user_id` (`user_id`),
    INDEX `idx_goal_feature_id` (`feature_id`),
    INDEX `idx_goal_variation_id` (`variation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 