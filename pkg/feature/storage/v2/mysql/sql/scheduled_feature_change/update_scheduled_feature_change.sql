UPDATE scheduled_feature_change
SET scheduled_at = ?,
    timezone = ?,
    payload = ?,
    comment = ?,
    status = ?,
    failure_reason = ?,
    conflicts = ?,
    updated_by = ?,
    updated_at = ?,
    executed_at = ?
WHERE id = ?
  AND environment_id = ?
