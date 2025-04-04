SELECT
    environment_v2.*,
    COALESCE(COUNT(DISTINCT feature.id), 0) AS feature_count
FROM
    environment_v2
        LEFT JOIN
    feature ON environment_v2.id = feature.environment_id
%s
GROUP BY
    environment_v2.id