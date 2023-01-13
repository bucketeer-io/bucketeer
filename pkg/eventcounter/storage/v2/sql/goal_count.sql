WITH grouped_by_user_evaluation AS (
    SELECT
        user_id,
        variation_id,
        COUNT(id) as event_count,
        SUM(value) as value_sum
    FROM
        `%s`
    WHERE
        _PARTITIONTIME BETWEEN TIMESTAMP(@startAt) AND TIMESTAMP(@endAt)
    AND environment_namespace = @environmentNamespace
    AND feature_id = @featureID
    AND feature_version = @featureVersion
GROUP BY
    user_id,
    variation_id
)
SELECT
    variation_id as variation,
    COUNT(DISTINCT user_id) as goal_user,
    SUM(event_count) as goal_total,
    SUM(value_sum) as goal_value_total,
    AVG(value_sum) as goal_value_mean,
    VAR_SAMP(value_sum) as goal_value_variance
FROM
    grouped_by_user_evaluation
GROUP BY
    variation_id
ORDER BY
    variation_id
