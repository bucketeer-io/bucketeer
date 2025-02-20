SELECT
    aor.id,
    aor.feature_id,
    aor.ops_type,
    aor.clauses,
    aor.created_at,
    aor.updated_at,
    aor.deleted,
    aor.status,
    ft.name
FROM
    auto_ops_rule as aor
JOIN feature ft ON
    aor.feature_id = ft.id AND
    aor.environment_id = ft.environment_id
%s %s %s
