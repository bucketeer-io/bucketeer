-- Universal holdout groups (RFC 0048, phase 4).
-- holdout_group_segment_user: the included users of the segment of a SEGMENT
-- group, copied when the group starts. user_id is compared byte for byte, the
-- way evaluation compares it, so it does not take the collation of the table.

CREATE TABLE IF NOT EXISTS holdout_group_segment_user (
    holdout_group_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    PRIMARY KEY (holdout_group_id, environment_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
