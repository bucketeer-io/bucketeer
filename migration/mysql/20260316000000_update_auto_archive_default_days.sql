-- Change the default value of auto_archive_unused_days from 90 to 60
ALTER TABLE `environment_v2` ALTER COLUMN `auto_archive_unused_days` SET DEFAULT 60;
