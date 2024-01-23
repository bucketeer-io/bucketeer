INSERT INTO auto_ops_rule (
    id,
    feature_id,
    ops_type,
    clauses,
    triggered_at,
    created_at,
    updated_at,
    deleted,
    environment_namespace
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
)
