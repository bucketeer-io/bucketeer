INSERT INTO monthly_summary (
    environment_id,
    source_id,
    yearmonth,
    mau,
    request_count,
    created_at,
    updated_at
) VALUES %s
ON DUPLICATE KEY UPDATE
    mau = VALUES(mau),
    updated_at = VALUES(updated_at)
