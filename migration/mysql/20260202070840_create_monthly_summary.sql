CREATE TABLE `monthly_summary` (
  `environment_id` varchar(255) NOT NULL,
  `source_id` varchar(30) NOT NULL,
  `yearmonth` varchar(6) NOT NULL,
  `mau` bigint NOT NULL DEFAULT 0,
  `request_count` bigint NOT NULL DEFAULT 0,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`environment_id`, `yearmonth`, `source_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
