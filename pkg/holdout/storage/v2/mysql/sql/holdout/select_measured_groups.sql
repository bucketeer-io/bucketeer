SELECT
    id,
    environment_id,
    name,
    description,
    assignment_type,
    sampling_seed,
    bucket_start,
    percentage,
    segment_id,
    status,
    start_at,
    stop_at,
    start_analysis_at,
    analysis_started_at,
    stopped_at,
    goal_ids,
    experiment_id,
    analysis_experiment_id,
    created_by,
    updated_by,
    created_at,
    updated_at,
    archived,
    deleted
FROM
    holdout_group
WHERE
    environment_id = ?
    AND deleted = ?
    AND (experiment_id <> '' OR analysis_experiment_id <> '')
    AND (
        status IN (?, ?)
        OR (status = ? AND stop_at >= ?)
        OR (status = ? AND stopped_at >= ?)
    )
ORDER BY
    start_at ASC,
    id ASC
