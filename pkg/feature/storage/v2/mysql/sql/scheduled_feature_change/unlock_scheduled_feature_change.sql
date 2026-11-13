UPDATE scheduled_feature_change
SET locked_at = NULL,
    locked_by = NULL
WHERE id = ?
  AND environment_id = ?
  AND locked_by = ?
