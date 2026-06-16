SELECT
    ms.environment_id,
    e.name AS environment_name,
    p.name AS project_name,
    ms.source_id,
    ms.yearmonth,
    ms.mau,
    ms.request_count
FROM monthly_summary ms
JOIN environment_v2 e ON ms.environment_id = e.id
JOIN project p ON e.project_id = p.id
WHERE ms.environment_id IN (%s) AND ms.source_id IN (%s)
ORDER BY ms.yearmonth ASC
