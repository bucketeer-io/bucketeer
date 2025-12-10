SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as evaluationUser,
    COUNT(id) as evaluationTotal
FROM
    evaluation_event
WHERE
    timestamp BETWEEN $1 AND $2
    AND environment_id = $3
    AND feature_id = $4
    AND feature_version = $5
GROUP BY
    variation_id 