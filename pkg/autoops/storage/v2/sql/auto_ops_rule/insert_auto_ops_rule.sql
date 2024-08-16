INSERT INTO auto_ops_rule (
    id,
    feature_id,
    ops_type,
    clauses,
    created_at,
    updated_at,
    deleted,
    status,
    environment_namespace
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
)
