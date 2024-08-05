UPDATE 
    auto_ops_rule
SET
    feature_id = ?,
    ops_type = ?,
    clauses = ?,
    triggered_at = 0,
    created_at = ?,
    updated_at = ?,
    deleted = ?,
    status = ?
WHERE
    id = ? AND
    environment_namespace = ?
