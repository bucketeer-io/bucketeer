UPDATE 
    ops_progressive_rollout
SET
    feature_id = ?,
    clause = ?,
    status = ?,
    type = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_namespace = ?
