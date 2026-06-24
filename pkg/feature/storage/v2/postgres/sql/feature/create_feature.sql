INSERT INTO feature (
    id,
    name,
    description,
    enabled,
    archived,
    deleted,
    evaluation_undelayable,
    ttl,
    version,
    created_at,
    updated_at,
    variation_type,
    variation_value_schema,
    variations,
    targets,
    rules,
    default_strategy,
    off_variation,
    tags,
    maintainer,
    sampling_seed,
    prerequisites,
    environment_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23
)
