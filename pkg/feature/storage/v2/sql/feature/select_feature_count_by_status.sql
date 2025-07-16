SELECT
    COUNT(DISTINCT feature.id) AS total_features,
    COUNT(DISTINCT CASE WHEN FROM_UNIXTIME(feature_last_used_info.last_used_at) >= DATE_SUB(now(), INTERVAL 7 DAY) THEN feature.id END) AS active,
    COUNT(DISTINCT CASE WHEN FROM_UNIXTIME(feature_last_used_info.last_used_at) < DATE_SUB(now(), INTERVAL 7 DAY) THEN feature.id END) AS inactive
FROM
    feature
LEFT JOIN feature_last_used_info ON
    feature.id = feature_last_used_info.feature_id AND
    feature.environment_id = feature_last_used_info.environment_id AND
    NOT EXISTS (
        SELECT 1 FROM feature_last_used_info flui2
        WHERE flui2.feature_id = feature_last_used_info.feature_id
        AND flui2.environment_id = feature_last_used_info.environment_id
        AND (flui2.last_used_at > feature_last_used_info.last_used_at
        OR (flui2.last_used_at = feature_last_used_info.last_used_at AND flui2.version > feature_last_used_info.version))
    )
WHERE
    feature.archived = 0 AND
    feature.deleted = 0 AND
    feature.environment_id = ?