-- Universal holdout groups (RFC 0048, phase 4).
-- kind tells an experiment on a feature flag (0) from an internal experiment
-- that measures a holdout group (1). Every existing row is an experiment on a
-- feature flag.
ALTER TABLE experiment ADD COLUMN kind INTEGER NOT NULL DEFAULT 0;

-- An internal experiment stores the holdout group ID in feature_id, and that
-- ID is not a row in feature. MySQL has had no such key since
-- 20241113142932_update_environment_id_table.sql. The index
-- idx_experiment_feature stays.
ALTER TABLE experiment DROP CONSTRAINT foreign_experiment_feature;
