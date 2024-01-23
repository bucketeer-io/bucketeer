INSERT INTO ops_progressive_rollout (
    id,
    feature_id,
    clause,
    status,
    type,
    created_at,
    updated_at,
    environment_namespace
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
