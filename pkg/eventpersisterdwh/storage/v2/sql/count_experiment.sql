SELECT
    (
        SELECT
            count(*)
        FROM
            experiment
        WHERE
            status = 1
    ) + (
        SELECT
            count(*)
        FROM
            experiment
        WHERE
            status = 2
          AND stop_at > UNIX_TIMESTAMP(
                DATE_SUB(NOW(), INTERVAL 2 DAY)
          )
    ) as count;