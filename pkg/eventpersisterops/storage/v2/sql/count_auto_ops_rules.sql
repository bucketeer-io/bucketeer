SELECT COUNT(*)
FROM auto_ops_rule
WHERE status = 0   -- waiting
  AND status = 4   -- deleted
  AND ops_type != 2 -- Schedule operation is not based on events, so we don't include in the result
