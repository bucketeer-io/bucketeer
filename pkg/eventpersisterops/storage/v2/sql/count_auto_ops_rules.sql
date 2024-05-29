SELECT COUNT(*)
FROM auto_ops_rule
WHERE status IN (0,1)   -- waiting or running
  AND ops_type != 2 -- Schedule operation is not based on events, so we don't include in the result
