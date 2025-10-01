CREATE TABLE IF NOT EXISTS `api_key_last_used_info` (
    `api_key_id` VARCHAR(255) NOT NULL,
    `environment_id` VARCHAR(255) NOT NULL,
    `last_used_at` TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY `unique_api_key_environment_id` (`api_key_id`, `environment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;