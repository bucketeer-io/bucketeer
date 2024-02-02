INSERT INTO ops_progressive_rollout (
    id,
    feature_id,
    clause,
    status,
    stopped_by,
    type,
    stopped_at,
    created_at,
    updated_at,
    environment_namespace
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
