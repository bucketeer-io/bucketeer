UPDATE
    experiment
SET
    goal_id = ?,
    feature_id = ?,
    feature_version = ?,
    variations = ?,
    start_at = ?,
    stop_at = ?,
    stopped = ?,
    stopped_at = ?,
    created_at = ?,
    updated_at = ?,
    archived = ?,
    deleted = ?,
    goal_ids = ?,
    name = ?,
    description = ?,
    base_variation_id = ?,
    maintainer = ?,
    status = ?
WHERE
    id = ? AND
    environment_id = ?