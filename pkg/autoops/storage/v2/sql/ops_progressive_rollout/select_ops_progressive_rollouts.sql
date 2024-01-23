SELECT
    id,
    feature_id,
    clause,
    status,
    type,
    created_at,
    updated_at
FROM
    ops_progressive_rollout
%s %s %s
