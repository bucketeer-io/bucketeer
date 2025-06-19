CREATE TABLE `team` (
    `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text,
    `organization_id` varchar(255) NOT NULL,
    `created_at` bigint NOT NULL,
    `updated_at` bigint NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_name_org` (`name`, `organization_id`),
    KEY `idx_organization_id` (`organization_id`)
);