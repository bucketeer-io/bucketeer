UPDATE scheduled_feature_change
SET locked_at = $1,
    locked_by = $2
WHERE id = $3
  AND (locked_at IS NULL OR locked_at < $4)
  AND status = 1
