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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
