SELECT
    COALESCE(SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END), 0) AS waiting,
    COALESCE(SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END), 0) AS running,
    COALESCE(SUM(CASE WHEN status = 2 OR status = 3 THEN 1 ELSE 0 END), 0) AS stopped
FROM
    experiment
WHERE
    environment_id = $1
    AND deleted = false
    AND archived = false
