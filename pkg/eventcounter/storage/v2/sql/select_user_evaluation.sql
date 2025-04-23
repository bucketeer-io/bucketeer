SELECT
    user_id as userID,
    feature_id as featureID,
    feature_version,
    variation_id as variationID,
    reason,
    UNIX_SECONDS(timestamp) as timestamp
FROM
    `%s`
WHERE
    environment_id = @environmentId
    AND feature_id = @featureId
    AND feature_version = @featureVersion
    AND user_id = @userId
    AND timestamp BETWEEN TIMESTAMP(@experimentStartAt)
    AND TIMESTAMP(@experimentEndAt)
ORDER BY
    timestamp DESC
LIMIT 1;
