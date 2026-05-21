SELECT
    id,
    segment_id,
    user_id,
    state,
    deleted
FROM segment_user
WHERE segment_id = $1
    AND environment_id = $2
    AND deleted = FALSE
