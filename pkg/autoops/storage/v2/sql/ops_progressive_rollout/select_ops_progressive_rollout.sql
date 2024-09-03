SELECT
    id,
    feature_id,
    clause,
    status,
    stopped_by,
    type,
    stopped_at,
    created_at,
    updated_at
FROM
    ops_progressive_rollout
WHERE
    id = ? AND
    environment_id = ?
