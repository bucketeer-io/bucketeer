-- Migrate admin_audit_log data to audit_log table
INSERT INTO `audit_log` (
    `id`,
    `timestamp`,
    `entity_type`,
    `entity_id`,
    `type`,
    `event`,
    `editor`,
    `options`,
    `entity_data`,
    `previous_entity_data`,
    `environment_id`
) SELECT 
    `id`,
    `timestamp`,
    `entity_type`,
    `entity_id`,
    `type`,
    `event`,
    `editor`,
    `options`,
    `entity_data`,
    `previous_entity_data`,
    '' as `environment_id`
FROM `admin_audit_log`;
