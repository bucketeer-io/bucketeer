UPDATE api_key SET
    name = $1,
    role = $2,
    disabled = $3,
    maintainer = $4,
    description = $5,
    updated_at = $6
WHERE
    id = $7 AND
    environment_id = $8
