UPDATE 
    auto_ops_rule
SET
    feature_id = ?,
    ops_type = ?,
    clauses = ?,
    triggered_at = ?,
    created_at = ?,
    updated_at = ?,
    deleted = ?
WHERE
    id = ? AND
    environment_namespace = ?
