SELECT
    id,
    feature_id,
    ops_type,
    clauses,
    triggered_at,
    created_at,
    updated_at,
    deleted,
    status
FROM
    auto_ops_rule
%s %s %s
