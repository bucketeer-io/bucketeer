SELECT
    COUNT(*)
FROM (
    SELECT
        1 AS included
    FROM
        segment_user
    WHERE
        segment_id = ? AND
        environment_id = ? AND
        state = ? AND
        deleted = 0
    LIMIT ?
) AS included_users
