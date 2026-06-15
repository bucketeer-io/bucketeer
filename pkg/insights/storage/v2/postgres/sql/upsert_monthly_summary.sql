INSERT INTO monthly_summary (
    environment_id,
    source_id,
    yearmonth,
    mau,
    request_count,
    created_at,
    updated_at
) VALUES %s
ON CONFLICT (environment_id, yearmonth, source_id) DO UPDATE SET
    mau = EXCLUDED.mau,
    request_count = CASE
        WHEN EXCLUDED.request_count = 0 THEN monthly_summary.request_count
        ELSE EXCLUDED.request_count
    END,
    updated_at = EXCLUDED.updated_at
