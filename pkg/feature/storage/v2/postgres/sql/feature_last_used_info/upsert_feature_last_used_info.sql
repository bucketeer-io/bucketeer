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
    $1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (id, environment_id) DO UPDATE SET
    feature_id = EXCLUDED.feature_id,
    version = EXCLUDED.version,
    last_used_at = EXCLUDED.last_used_at,
    client_oldest_version = EXCLUDED.client_oldest_version,
    client_latest_version = EXCLUDED.client_latest_version