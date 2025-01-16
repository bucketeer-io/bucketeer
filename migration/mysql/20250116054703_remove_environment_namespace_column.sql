-- Remove environment_namespace column from all tables
ALTER TABLE feature DROP COLUMN environment_namespace;
ALTER TABLE api_key DROP COLUMN environment_namespace;
ALTER TABLE audit_log DROP COLUMN environment_namespace;
ALTER TABLE auto_ops_rule DROP COLUMN environment_namespace;
ALTER TABLE experiment DROP COLUMN environment_namespace;
ALTER TABLE experiment_result DROP COLUMN environment_namespace;
ALTER TABLE feature_last_used_info DROP COLUMN environment_namespace;
ALTER TABLE flag_trigger DROP COLUMN environment_namespace;
ALTER TABLE goal DROP COLUMN environment_namespace;
-- We don't use mau table anymore, ignore it
-- ALTER TABLE mau DROP COLUMN environment_namespace;
ALTER TABLE ops_count DROP COLUMN environment_namespace;
ALTER TABLE ops_progressive_rollout DROP COLUMN environment_namespace;
ALTER TABLE push DROP COLUMN environment_namespace;
ALTER TABLE segment DROP COLUMN environment_namespace;
ALTER TABLE segment_user DROP COLUMN environment_namespace;
ALTER TABLE subscription DROP COLUMN environment_namespace;
ALTER TABLE tag DROP COLUMN environment_namespace;