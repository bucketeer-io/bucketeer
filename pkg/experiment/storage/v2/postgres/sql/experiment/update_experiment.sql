UPDATE
    experiment
SET
    goal_id = $1,
    feature_id = $2,
    feature_version = $3,
    variations = $4,
    start_at = $5,
    stop_at = $6,
    stopped = $7,
    stopped_at = $8,
    created_at = $9,
    updated_at = $10,
    archived = $11,
    deleted = $12,
    goal_ids = $13,
    name = $14,
    description = $15,
    base_variation_id = $16,
    maintainer = $17,
    status = $18
WHERE
    id = $19 AND
    environment_id = $20
