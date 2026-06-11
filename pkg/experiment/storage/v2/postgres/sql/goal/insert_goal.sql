INSERT INTO goal (
    id,
    name,
    description,
    connection_type,
    archived,
    deleted,
    created_at,
    updated_at,
    environment_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
