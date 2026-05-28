UPDATE api_key
SET last_used_at = $1
WHERE id = $2
  AND environment_id = $3
  AND last_used_at < $4
