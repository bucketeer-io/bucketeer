SELECT COUNT(*)
FROM auto_ops_rule
LEFT JOIN feature
ON auto_ops_rule.feature_id = feature.id
AND auto_ops_rule.environment_namespace = feature.environment_namespace
WHERE feature.archived = 0
  AND auto_ops_rule.triggered_at = 0
  AND auto_ops_rule.deleted = 0
  AND JSON_EXTRACT(JSON_EXTRACT(auto_ops_rule.clauses, '$[*].clause'), '$[0].type_url') != 'type.googleapis.com/bucketeer.autoops.DatetimeClause'
  -- Schedule operation is not based on events, so we don't include in the result
