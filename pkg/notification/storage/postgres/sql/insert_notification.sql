INSERT INTO notification (
    id,
    status,
    created_by,
    last_edited_by,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
