INSERT INTO api_key_last_used_info (
    api_key_id,
    last_used_at,
    created_at,
    environment_id
) VALUES (
    ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
    last_used_at = VALUES(last_used_at)