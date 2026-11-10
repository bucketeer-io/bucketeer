SELECT
    COUNT(DISTINCT feature.id) AS total_features,
    COUNT(DISTINCT CASE WHEN
        feature.last_used_at >= UNIX_TIMESTAMP(DATE_SUB(now(), INTERVAL 7 DAY))
        THEN feature.id
    END) AS active,
    COUNT(DISTINCT CASE WHEN
        feature.last_used_at > 0 AND
        feature.last_used_at < UNIX_TIMESTAMP(DATE_SUB(now(), INTERVAL 7 DAY))
        THEN feature.id
    END) AS inactive
FROM
    feature
WHERE
    feature.archived = 0 AND
    feature.deleted = 0 AND
    feature.environment_id = ?