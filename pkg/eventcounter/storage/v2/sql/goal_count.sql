WITH grouped_by_user_evaluation AS (
    SELECT
        user_id,
        variation_id,
        COUNT(id) as event_count,
        IFNULL(SUM(value), 0) as value_sum
    FROM
        `%s`
    WHERE
        timestamp BETWEEN TIMESTAMP(@startAt) AND TIMESTAMP(@endAt)
    AND environment_namespace = @environmentNamespace
    AND goal_id = @goalID
    AND feature_id = @featureID
    AND feature_version = @featureVersion
GROUP BY
    user_id,
    variation_id
)
SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as goalUser,
    SUM(event_count) as goalTotal,
    SUM(value_sum) as goalValueTotal,
    AVG(value_sum) as goalValueMean,
    IFNULL(VAR_SAMP(value_sum), 0) as goalValueVariance
FROM
    grouped_by_user_evaluation
GROUP BY
    variation_id
ORDER BY
    variation_id
