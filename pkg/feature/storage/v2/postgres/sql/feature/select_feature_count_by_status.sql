SELECT
    COUNT(DISTINCT feature.id) AS total_features,
    COUNT(DISTINCT CASE WHEN
        feature.last_used_at >= EXTRACT(EPOCH FROM NOW())::BIGINT - 604800
        THEN feature.id
    END) AS active,
    COUNT(DISTINCT CASE WHEN
        feature.last_used_at > 0 AND
        feature.last_used_at < EXTRACT(EPOCH FROM NOW())::BIGINT - 604800
        THEN feature.id
    END) AS inactive
FROM
    feature
WHERE
    feature.archived = FALSE AND
    feature.deleted = FALSE AND
    feature.environment_id = $1
