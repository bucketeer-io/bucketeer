-- Create "account_credentials" table for password authentication
CREATE TABLE `account_credentials` (
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `password_reset_token` varchar(255) DEFAULT NULL,
  `password_reset_token_expires_at` bigint DEFAULT NULL,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`email`),
  UNIQUE INDEX `unique_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- Create index for password reset tokens to improve lookup performance
CREATE INDEX `idx_password_reset_token` ON `account_credentials` (`password_reset_token`);