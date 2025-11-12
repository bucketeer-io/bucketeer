SELECT
    (
        SELECT
            count(*)
        FROM
            experiment
        WHERE
            status = 1
        AND archived = 0
    ) + (
        SELECT
            count(*)
        FROM
            experiment
        WHERE
            status = 2
        AND archived = 0
        AND stop_at > UNIX_TIMESTAMP(
            DATE_SUB(NOW(), INTERVAL 2 DAY)
        )
    ) as count;