UPDATE scheduled_feature_change
SET locked_at = ?,
    locked_by = ?
WHERE id = ?
  AND (locked_at IS NULL OR locked_at < ?)
  AND status = 1
