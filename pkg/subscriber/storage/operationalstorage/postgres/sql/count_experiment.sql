SELECT
    (
        SELECT
            count(*)
        FROM
            experiment
        WHERE
            status = 1
        AND archived = false
    ) + (
        SELECT
            count(*)
        FROM
            experiment
        WHERE
            status = 2
        AND archived = false
        AND stop_at > EXTRACT(EPOCH FROM (NOW() - INTERVAL '2 days'))::bigint
    ) as count;
