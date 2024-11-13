-- Step 1: Drop all foreign keys referencing environment_namespace
ALTER TABLE auto_ops_rule DROP FOREIGN KEY foreign_auto_ops_rule_feature_id_environment_namespace;
ALTER TABLE experiment DROP FOREIGN KEY foreign_experiment_feature_id_environment_namespace;
ALTER TABLE ops_count 
DROP FOREIGN KEY foreign_ops_count_feature_id_environment_namespace,
DROP FOREIGN KEY foreign_ops_count_auto_ops_rule_id_environment_namespace;
ALTER TABLE ops_progressive_rollout DROP FOREIGN KEY foreign_progressive_rollout_feature_id_environment_namespace;
ALTER TABLE flag_trigger DROP FOREIGN KEY foreign_flag_trigger_feature_id_environment_namespace;
ALTER TABLE segment_user DROP FOREIGN KEY foreign_segment_user_segment_id_environment_namespace;

-- Step 2: Populate environment_id with values from environment_namespace
UPDATE feature SET environment_id = environment_namespace;
UPDATE account SET environment_id = environment_namespace;
UPDATE api_key SET environment_id = environment_namespace;
UPDATE audit_log SET environment_id = environment_namespace;
UPDATE auto_ops_rule SET environment_id = environment_namespace;
UPDATE experiment SET environment_id = environment_namespace;
UPDATE experiment_result SET environment_id = environment_namespace;
UPDATE feature_last_used_info SET environment_id = environment_namespace;
UPDATE flag_trigger SET environment_id = environment_namespace;
UPDATE goal SET environment_id = environment_namespace;
UPDATE mau SET environment_id = environment_namespace;
UPDATE ops_count SET environment_id = environment_namespace;
UPDATE ops_progressive_rollout SET environment_id = environment_namespace;
UPDATE push SET environment_id = environment_namespace;
UPDATE segment SET environment_id = environment_namespace;
UPDATE segment_user SET environment_id = environment_namespace;
UPDATE subscription SET environment_id = environment_namespace;
UPDATE tag SET environment_id = environment_namespace;

-- Step 3: Change PRIMARY KEY to use environment_id instead of environment_namespace
ALTER TABLE feature_last_used_info
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE mau
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (user_id, yearmonth, source_id, environment_id);

ALTER TABLE ops_count
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE ops_progressive_rollout
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE segment_user
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE feature
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE account
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE api_key
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE audit_log
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE auto_ops_rule
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE experiment
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE experiment_result
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE flag_trigger
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE goal
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE push
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE segment
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE subscription
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);

ALTER TABLE tag
ALGORITHM=INPLACE,
DROP PRIMARY KEY,
ADD PRIMARY KEY (id, environment_id);