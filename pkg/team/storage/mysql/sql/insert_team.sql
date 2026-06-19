INSERT INTO team (
    id,
    name,
    description,
    organization_id,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
    updated_at = VALUES(updated_at)
