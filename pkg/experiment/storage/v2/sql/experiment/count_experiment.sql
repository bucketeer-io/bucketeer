SELECT
    COUNT(1) AS total,
    COALESCE(SUM(IF(status = 0, 1, 0)), 0) AS waiting,
    COALESCE(SUM(IF(status = 1, 1, 0)), 0) AS running,
    COALESCE(SUM(IF(status = 2 OR status = 3, 1, 0)), 0) AS stopped
FROM
    experiment
%s