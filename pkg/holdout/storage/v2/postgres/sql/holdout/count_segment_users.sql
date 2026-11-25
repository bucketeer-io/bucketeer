SELECT
    COUNT(*)
FROM (
    SELECT
        1 AS included
    FROM
        segment_user
    WHERE
        segment_id = $1 AND
        environment_id = $2 AND
        state = $3 AND
        deleted = FALSE
    LIMIT $4
) AS included_users
