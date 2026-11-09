SELECT
    COUNT(1)
FROM
    feature
LEFT JOIN LATERAL (
    SELECT e.id
    FROM experiment e
    WHERE e.feature_id = feature.id
    AND e.environment_id = feature.environment_id
    AND e.deleted = FALSE
    LIMIT 1
) experiment ON TRUE
