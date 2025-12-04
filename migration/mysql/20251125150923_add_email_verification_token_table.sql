-- Create "email_verification_token" table for magic link authentication
CREATE TABLE `email_verification_token` (
  `email` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `created_at` bigint NOT NULL,
  `expires_at` bigint NOT NULL,
  `verified_at` bigint DEFAULT NULL,
  `ip_address` varchar(45) DEFAULT NULL,
  `user_agent` text DEFAULT NULL,
  PRIMARY KEY (`email`),
  UNIQUE INDEX `idx_token` (`token`),
  INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
