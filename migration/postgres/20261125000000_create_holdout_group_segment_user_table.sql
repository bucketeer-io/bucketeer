-- Universal holdout groups (RFC 0048, phase 4).
-- holdout_group_segment_user: the included users of the segment of a SEGMENT
-- group, copied when the group starts.

CREATE TABLE holdout_group_segment_user (
    holdout_group_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (holdout_group_id, environment_id, user_id)
);
