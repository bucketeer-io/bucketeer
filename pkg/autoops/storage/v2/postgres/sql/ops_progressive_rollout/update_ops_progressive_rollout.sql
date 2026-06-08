UPDATE
    ops_progressive_rollout
SET
    feature_id = $1,
    clause = $2,
    status = $3,
    stopped_by = $4,
    type = $5,
    stopped_at = $6,
    created_at = $7,
    updated_at = $8
WHERE
    id = $9 AND
    environment_id = $10
