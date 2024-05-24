SELECT COUNT(*)
FROM auto_ops_rule
WHERE triggered_at = 0
  AND deleted = 0
  AND JSON_EXTRACT(JSON_EXTRACT(clauses, '$[*].clause'), '$[0].type_url') != 'type.googleapis.com/bucketeer.autoops.DatetimeClause'
  -- Schedule operation is not based on events, so we don't include in the result
