SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as evaluationUser,
    COUNT(id) as evaluationTotal
FROM
    evaluation_event
WHERE
    timestamp BETWEEN ? AND ?
    AND environment_id = ?
    AND feature_id = ?
    AND feature_version = ?
GROUP BY
    variation_id 