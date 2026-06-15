INSERT INTO experiment (
    id,
    goal_id,
    feature_id,
    feature_version,
    variations,
    start_at,
    stop_at,
    stopped,
    stopped_at,
    created_at,
    updated_at,
    archived,
    deleted,
    goal_ids,
    name,
    description,
    base_variation_id,
    status,
    maintainer,
    environment_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
