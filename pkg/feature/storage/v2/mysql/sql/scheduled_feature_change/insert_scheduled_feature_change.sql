INSERT INTO scheduled_feature_change (
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
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
