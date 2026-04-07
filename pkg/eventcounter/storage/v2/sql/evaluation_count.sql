SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as evaluationUser,
    COUNT(id) as evaluationTotal
FROM
    `%s`
WHERE
    timestamp BETWEEN @startAt AND @endAt
    AND environment_id = @environmentId
    AND feature_id = @featureID
    AND feature_version = @featureVersion
GROUP BY
    variation_id
