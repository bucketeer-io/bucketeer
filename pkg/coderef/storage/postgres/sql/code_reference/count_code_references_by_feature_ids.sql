SELECT feature_id, COUNT(1) as count
FROM code_reference
WHERE environment_id = $1
  AND feature_id IN (%s)
GROUP BY feature_id
