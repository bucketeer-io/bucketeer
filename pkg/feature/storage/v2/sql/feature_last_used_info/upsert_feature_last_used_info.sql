INSERT INTO feature_last_used_info (
    id,
    feature_id,
    version,
    last_used_at,
    client_oldest_version,
    client_latest_version,
    created_at,
    environment_id
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
    feature_id = VALUES(feature_id),
    version = VALUES(version),
    last_used_at = VALUES(last_used_at),
    client_oldest_version = VALUES(client_oldest_version),
    client_latest_version = VALUES(client_latest_version)