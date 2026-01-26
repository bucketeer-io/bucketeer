UPDATE scheduled_feature_change
SET locked_at = NULL,
    locked_by = NULL
WHERE id = ?
