INSERT INTO holdout_group_segment_user (
    holdout_group_id,
    environment_id,
    user_id
)
SELECT
    $1::VARCHAR,
    environment_id,
    user_id
FROM
    segment_user
WHERE
    segment_id = $2 AND
    environment_id = $3 AND
    state = $4 AND
    deleted = FALSE
