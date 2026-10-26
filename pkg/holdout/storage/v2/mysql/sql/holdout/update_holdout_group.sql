UPDATE holdout_group
SET
    name = ?,
    description = ?,
    status = ?,
    start_at = ?,
    stop_at = ?,
    start_analysis_at = ?,
    analysis_started_at = ?,
    stopped_at = ?,
    goal_ids = ?,
    experiment_id = ?,
    analysis_experiment_id = ?,
    updated_by = ?,
    updated_at = ?,
    archived = ?,
    deleted = ?
WHERE
    id = ? AND environment_id = ? AND deleted = ?
