-- Create "domain_auth_policy" table for domain-based authentication policies
CREATE TABLE `domain_auth_policy` (
  `domain` varchar(255) NOT NULL,
  `auth_policy` json NOT NULL,
  `enabled` boolean NOT NULL DEFAULT TRUE,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`domain`),
  UNIQUE INDEX `unique_domain` (`domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- Create index on enabled column to improve queries for active policies
CREATE INDEX `idx_enabled` ON `domain_auth_policy` (`enabled`);

-- Create index on created_at for ordering
CREATE INDEX `idx_created_at` ON `domain_auth_policy` (`created_at`);

-- Create index on updated_at for ordering
CREATE INDEX `idx_updated_at` ON `domain_auth_policy` (`updated_at`);
