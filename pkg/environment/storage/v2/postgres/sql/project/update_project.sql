UPDATE project
SET
    name = $1,
    description = $2,
    disabled = $3,
    trial = $4,
    creator_email = $5,
    created_at = $6,
    updated_at = $7
WHERE
    id = $8
