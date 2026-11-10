UPDATE scheduled_feature_change
SET locked_at = NULL,
    locked_by = NULL
WHERE id = $1
  AND environment_id = $2
  AND locked_by = $3
