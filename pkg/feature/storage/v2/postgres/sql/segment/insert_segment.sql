INSERT INTO segment (
    id,
    name,
    description,
    rules,
    created_at,
    updated_at,
    version,
    deleted,
    included_user_count,
    excluded_user_count,
    status,
    environment_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
