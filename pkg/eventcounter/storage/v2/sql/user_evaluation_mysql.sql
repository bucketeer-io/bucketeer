SELECT
    user_id as userID,
    feature_id as featureID,
    feature_version as featureVersion,
    variation_id as variationID,
    reason,
    CAST(UNIX_TIMESTAMP(timestamp) AS SIGNED) as timestamp
FROM
    evaluation_event
WHERE
    environment_id = ?
    AND feature_id = ?
    AND feature_version = ?
    AND user_id = ?
    AND timestamp BETWEEN ? AND ?
ORDER BY
    timestamp DESC
LIMIT 1; 