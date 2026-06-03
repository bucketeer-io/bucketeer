INSERT INTO auto_ops_rule (
    id,
    feature_id,
    ops_type,
    clauses,
    created_at,
    updated_at,
    deleted,
    status,
    environment_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
