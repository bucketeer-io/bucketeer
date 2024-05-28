SELECT
    id,
    feature_id,
    ops_type,
    clauses,
    created_at,
    updated_at,
    stopped_at,
    status
FROM
    auto_ops_rule
WHERE
    id = ? AND
    environment_namespace = ?
