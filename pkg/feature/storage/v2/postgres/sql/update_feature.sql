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
    variations = $12,
    targets = $13,
    rules = $14,
    default_strategy = $15,
    off_variation = $16,
    tags = $17,
    maintainer = $18,
    sampling_seed = $19,
    prerequisites = $20
WHERE
    id = $21 AND
    environment_id = $22
