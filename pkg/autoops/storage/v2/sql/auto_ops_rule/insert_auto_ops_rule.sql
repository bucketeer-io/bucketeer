INSERT INTO auto_ops_rule (
    id,
    feature_id,
    ops_type,
    clauses,
    created_at,
    updated_at,
    stopped_at,
    status,
    environment_namespace
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
)
