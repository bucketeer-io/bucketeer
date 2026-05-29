UPDATE
    organization
SET
    name = $1,
    owner_email = $2,
    description = $3,
    disabled = $4,
    archived = $5,
    trial = $6,
    created_at = $7,
    updated_at = $8
WHERE
    id = $9
