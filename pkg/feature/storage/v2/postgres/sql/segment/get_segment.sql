SELECT
    segment.id,
    segment.name,
    segment.description,
    segment.rules,
    segment.created_at,
    segment.updated_at,
    segment.version,
    segment.deleted,
    segment.included_user_count,
    segment.excluded_user_count,
    segment.status,
    (
        SELECT
            STRING_AGG(id, ',')
        FROM
            feature
        WHERE
            environment_id = $1 AND
            rules::text LIKE '%' || segment.id || '%'
    ) AS feature_ids
FROM
    segment
WHERE
    id = $2 AND
    environment_id = $3
