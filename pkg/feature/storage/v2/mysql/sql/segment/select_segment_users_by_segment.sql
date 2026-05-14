SELECT
    id,
    segment_id,
    user_id,
    state,
    deleted
FROM segment_user
WHERE segment_id = ?
    AND environment_id = ?
    AND deleted = 0;
