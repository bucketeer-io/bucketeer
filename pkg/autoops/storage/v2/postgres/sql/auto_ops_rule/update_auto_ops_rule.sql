UPDATE
    auto_ops_rule
SET
    feature_id = $1,
    ops_type = $2,
    clauses = $3,
    created_at = $4,
    updated_at = $5,
    deleted = $6,
    status = $7
WHERE
    id = $8 AND
    environment_id = $9
