UPDATE scheduled_feature_change
SET scheduled_at = $1,
    timezone = $2,
    payload = $3,
    comment = $4,
    status = $5,
    failure_reason = $6,
    conflicts = $7,
    updated_by = $8,
    updated_at = $9,
    executed_at = $10
WHERE id = $11
  AND environment_id = $12
