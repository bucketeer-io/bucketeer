UPDATE holdout_group
SET
    name = $1,
    description = $2,
    status = $3,
    start_at = $4,
    stop_at = $5,
    start_analysis_at = $6,
    analysis_started_at = $7,
    stopped_at = $8,
    goal_ids = $9,
    experiment_id = $10,
    analysis_experiment_id = $11,
    updated_by = $12,
    updated_at = $13,
    archived = $14,
    deleted = $15
WHERE
    id = $16 AND environment_id = $17 AND deleted = $18
