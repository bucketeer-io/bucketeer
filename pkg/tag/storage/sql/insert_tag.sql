INSERT INTO tag (
    id,
    created_at,
    updated_at,
    entity_type,
    environment_id
) VALUES (
    ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
    updated_at = VALUES(updated_at)
