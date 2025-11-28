UPDATE api_key
SET last_used_at = ?
WHERE id = ?
  AND environment_id = ?
  AND last_used_at < ?