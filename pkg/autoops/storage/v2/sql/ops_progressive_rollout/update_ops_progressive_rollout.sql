UPDATE 
    ops_progressive_rollout
SET
    feature_id = ?,
    clause = ?,
    status = ?,
    stopped_by = ?,
    type = ?,
    stopped_at = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_namespace = ?
