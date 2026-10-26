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
    environment_id = ?
    AND deleted = ?
    AND status IN (?, ?, ?)
