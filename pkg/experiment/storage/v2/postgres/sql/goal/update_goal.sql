UPDATE
    goal
SET
    name = $1,
    description = $2,
    archived = $3,
    deleted = $4,
    created_at = $5,
    updated_at = $6
WHERE
    id = $7 AND
    environment_id = $8
