SELECT
    user_id as userID,
    feature_id as featureID,
    feature_version as featureVersion,
    variation_id as variationID,
    reason,
    CAST(EXTRACT(EPOCH FROM "timestamp") AS BIGINT) as "timestamp"
FROM
    evaluation_event
WHERE
    environment_id = $1
    AND feature_id = $2
    AND feature_version = $3
    AND user_id = $4
    AND "timestamp" BETWEEN $5 AND $6
ORDER BY
    "timestamp" DESC
LIMIT 1; 