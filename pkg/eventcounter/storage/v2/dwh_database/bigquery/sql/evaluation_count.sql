SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as evaluationUser,
    -- DISTINCT by event ID: BigQuery has no primary keys, so an at-least-once
    -- Pub/Sub redelivery can append the same event twice. Each legitimate
    -- event has its own unique ID, so this still counts repeat evaluations
    -- by the same user.
    COUNT(DISTINCT id) as evaluationTotal
FROM
    `%s`
WHERE
    timestamp BETWEEN @startAt AND @endAt
    AND environment_id = @environmentId
    AND feature_id = @featureID
    AND feature_version = @featureVersion
GROUP BY
    variation_id
