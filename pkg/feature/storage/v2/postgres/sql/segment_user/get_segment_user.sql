SELECT
    id,
    segment_id,
    user_id,
    state,
    deleted
FROM
    segment_user
WHERE
    id = $1 AND
    environment_id = $2
