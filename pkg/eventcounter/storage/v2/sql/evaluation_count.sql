SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as evaluationUser,
    COUNT(id) as evaluationTotal
FROM
    `%s`
WHERE
    _PARTITIONTIME BETWEEN TIMESTAMP(@startAt) AND TIMESTAMP(@endAt)
    AND environment_namespace = @environmentNamespace
    AND feature_id = @featureID
    AND feature_version = @featureVersion
GROUP BY
    variation_id
