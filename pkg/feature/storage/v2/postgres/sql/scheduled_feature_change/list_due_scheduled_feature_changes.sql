SELECT
    id,
    feature_id,
    environment_id,
    scheduled_at,
    timezone,
    payload,
    comment,
    status,
    failure_reason,
    flag_version_at_creation,
    conflicts,
    created_by,
    created_at,
    updated_by,
    updated_at,
    executed_at
FROM scheduled_feature_change
WHERE scheduled_at <= $1
  AND status = 1
  AND (locked_at IS NULL OR locked_at < $2)
ORDER BY scheduled_at ASC
LIMIT $3
