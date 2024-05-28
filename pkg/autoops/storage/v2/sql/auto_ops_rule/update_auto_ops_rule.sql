UPDATE 
    auto_ops_rule
SET
    feature_id = ?,
    ops_type = ?,
    clauses = ?,
    created_at = ?,
    updated_at = ?,
    stopped_at = ?,
    status = ?
WHERE
    id = ? AND
    environment_namespace = ?
