-- Add account index
ALTER TABLE `account_v2` ADD INDEX `idx_email` (`email`);

-- Add organization indexes
ALTER TABLE `organization` ADD INDEX `idx_disabled` (`disabled`);
ALTER TABLE `organization` ADD INDEX `idx_archived` (`archived`);
ALTER TABLE `organization` ADD INDEX `idx_disabled_archived_id` (`disabled`, `archived`, `id`);

-- Add project index
ALTER TABLE `project` ADD INDEX `idx_organization_id` (`organization_id`);

-- Add environment_v2 index
ALTER TABLE `environment_v2` ADD INDEX `idx_organization_id` (`organization_id`);
