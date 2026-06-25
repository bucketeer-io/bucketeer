UPDATE
    feature
SET
    name = $1,
    description = $2,
    enabled = $3,
    archived = $4,
    deleted = $5,
    evaluation_undelayable = $6,
    ttl = $7,
    version = $8,
    created_at = $9,
    updated_at = $10,
    variation_type = $11,
    variation_value_schema = $12,
    variations = $13,
    targets = $14,
    rules = $15,
    default_strategy = $16,
    off_variation = $17,
    tags = $18,
    maintainer = $19,
    sampling_seed = $20,
    prerequisites = $21
WHERE
    id = $22 AND
    environment_id = $23
