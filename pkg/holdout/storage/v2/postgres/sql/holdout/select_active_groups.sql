SELECT
    id,
    name,
    status,
    start_at,
    stop_at,
    segment_id
FROM
    holdout_group
WHERE
    environment_id = $1
    AND deleted = $2
    AND status IN ($3, $4, $5)
