SELECT
    user_id as userID,
    feature_id as featureID,
    feature_version,
    variation_id as variationID,
    reason,
    timestamp
FROM
    `%s`
WHERE
    timestamp BETWEEN TIMESTAMP(@experimentStartAt)
    AND LEAST(
        TIMESTAMP(@experimentEndAt)
        TIMESTAMP(@goalTimestamp)
    )
    AND environment_id = @environmentId
    AND feature_id = @featureId
    AND user_id = @userId

ORDER BY
    timestamp DESC
LIMIT 1;
