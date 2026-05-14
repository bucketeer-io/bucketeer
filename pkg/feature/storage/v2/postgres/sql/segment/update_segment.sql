UPDATE
    segment
SET
    name = $1,
    description = $2,
    rules = $3,
    created_at = $4,
    updated_at = $5,
    version = $6,
    deleted = $7,
    included_user_count = $8,
    excluded_user_count = $9,
    status = $10
WHERE
    id = $11 AND
    environment_id = $12
