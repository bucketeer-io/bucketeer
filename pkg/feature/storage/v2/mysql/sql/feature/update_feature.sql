UPDATE
    feature
SET
    name = ?,
    description = ?,
    enabled = ?,
    archived = ?,
    deleted = ?,
    evaluation_undelayable = ?,
    ttl = ?,
    version = ?,
    created_at = ?,
    updated_at = ?,
    variation_type = ?,
    variation_value_schema = ?,
    variations = ?,
    targets = ?,
    rules = ?,
    default_strategy = ?,
    off_variation = ?,
    tags = ?,
    maintainer = ?,
    sampling_seed = ?,
    prerequisites = ?
WHERE
    id = ? AND
    environment_id = ?
