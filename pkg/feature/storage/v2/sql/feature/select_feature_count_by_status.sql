SELECT
    COUNT(DISTINCT feature.id) AS total_features,
    COUNT(DISTINCT CASE WHEN FROM_UNIXTIME(last_used_at) >= DATE_SUB(now(), INTERVAL 7 DAY) THEN feature.id END) AS active,
    COUNT(DISTINCT CASE WHEN FROM_UNIXTIME(last_used_at) < DATE_SUB(now(), INTERVAL 7 DAY) THEN feature.id END) AS inactive
FROM feature
    LEFT JOIN feature_last_used_info flui ON feature.id = flui.feature_id
WHERE
    feature.archived = 0 AND
    feature.deleted = 0 AND
    feature.environment_id = ?