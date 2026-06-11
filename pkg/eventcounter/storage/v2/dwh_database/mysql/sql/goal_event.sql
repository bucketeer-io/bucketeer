WITH grouped_by_user_evaluation AS (
    SELECT
        user_id,
        variation_id,
        COUNT(id) as event_count,
        IFNULL(SUM(value), 0) as value_sum
    FROM
        goal_event
    WHERE
        timestamp BETWEEN ? AND ?
        AND environment_id = ?
        AND goal_id = ?
        AND feature_id = ?
        AND feature_version = ?
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
    IFNULL(VARIANCE(value_sum), 0) as goalValueVariance
FROM
    grouped_by_user_evaluation
GROUP BY
    variation_id
ORDER BY
    variation_id 