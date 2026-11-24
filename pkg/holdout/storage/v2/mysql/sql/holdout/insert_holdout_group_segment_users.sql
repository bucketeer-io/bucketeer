INSERT INTO holdout_group_segment_user (
    holdout_group_id,
    environment_id,
    user_id
)
SELECT
    ?,
    environment_id,
    user_id
FROM
    segment_user
WHERE
    segment_id = ? AND
    environment_id = ? AND
    state = ? AND
    deleted = 0
