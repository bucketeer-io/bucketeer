SELECT COUNT(*)
FROM auto_ops_rule
LEFT JOIN feature
ON auto_ops_rule.feature_id = feature.id
AND auto_ops_rule.environment_namespace = feature.environment_namespace
WHERE feature.archived = 0
  AND auto_ops_rule.status IN (0, 1) -- 0: WAITING, 1: RUNNING
  AND auto_ops_rule.deleted = 0
  AND auto.ops_rule.ops_type != 2
  -- Schedule operation is not based on events, so we don't include in the result
