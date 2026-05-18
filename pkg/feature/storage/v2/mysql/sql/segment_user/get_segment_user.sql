SELECT
    id,
    segment_id,
    user_id,
    state,
    deleted
FROM
    segment_user
WHERE
    id = ? AND
    environment_id = ?
