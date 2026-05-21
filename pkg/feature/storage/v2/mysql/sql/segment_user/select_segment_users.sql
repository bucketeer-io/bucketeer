SELECT
    id,
    segment_id,
    user_id,
    state,
    deleted
FROM
    segment_user
%s %s
