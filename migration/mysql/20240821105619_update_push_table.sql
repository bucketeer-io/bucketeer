-- Modify "push" table
ALTER TABLE `push` ADD COLUMN `fcm_service_account` json NOT NULL;
