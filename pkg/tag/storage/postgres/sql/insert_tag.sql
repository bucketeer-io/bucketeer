INSERT INTO tag (
    id,
    name,
    created_at,
    updated_at,
    entity_type,
    environment_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT (name, environment_id, entity_type) DO UPDATE SET
    updated_at = EXCLUDED.updated_at
