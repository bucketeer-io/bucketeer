-- Step 1: Drop all foreign keys referencing environment_namespace
ALTER TABLE auto_ops_rule DROP FOREIGN KEY foreign_auto_ops_rule_feature_id_environment_namespace;
ALTER TABLE experiment DROP FOREIGN KEY foreign_experiment_feature_id_environment_namespace;
ALTER TABLE ops_count
DROP FOREIGN KEY foreign_ops_count_feature_id_environment_namespace,
  DROP FOREIGN KEY foreign_ops_count_auto_ops_rule_id_environment_namespace;
ALTER TABLE ops_progressive_rollout DROP FOREIGN KEY foreign_progressive_rollout_feature_id_environment_namespace;
ALTER TABLE flag_trigger DROP FOREIGN KEY foreign_flag_trigger_feature_id_environment_namespace;
ALTER TABLE segment_user DROP FOREIGN KEY foreign_segment_user_segment_id_environment_namespace;

-- Step 2: Change PRIMARY KEY to use environment_id instead of environment_namespace
-- and allow environment_namespace to be NULL

-- For tables with special consideration:
ALTER TABLE feature_last_used_info
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE mau
DROP PRIMARY KEY,
  ADD PRIMARY KEY (user_id, yearmonth, source_id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE ops_count
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE ops_progressive_rollout
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE segment_user
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

-- For other tables:
ALTER TABLE feature
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE account
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE api_key
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE audit_log
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE auto_ops_rule
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE experiment
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE experiment_result
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE flag_trigger
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE goal
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE push
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE segment
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE subscription
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

ALTER TABLE tag
DROP PRIMARY KEY,
  ADD PRIMARY KEY (id, environment_id),
  MODIFY COLUMN environment_namespace VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL;

-- Step 3: Populate environment_id with values from environment_namespace
UPDATE feature SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE account SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE api_key SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE audit_log SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE auto_ops_rule SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE experiment SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE experiment_result SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE feature_last_used_info SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE flag_trigger SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE goal SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE mau SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE ops_count SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE ops_progressive_rollout SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE push SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE segment SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE segment_user SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE subscription SET environment_id = environment_namespace WHERE environment_id IS NULL;
UPDATE tag SET environment_id = environment_namespace WHERE environment_id IS NULL;