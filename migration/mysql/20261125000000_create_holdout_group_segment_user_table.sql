-- holdout_group_segment_user: the included users of the segment of a SEGMENT
-- group, copied when the group starts. user_id is compared byte for byte, the
-- way evaluation compares it: case and trailing spaces count. utf8mb4_bin and
-- the collation of the table pad trailing spaces, so the column has
-- utf8mb4_0900_bin, which does not.

CREATE TABLE IF NOT EXISTS holdout_group_segment_user (
    holdout_group_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_bin NOT NULL,
    PRIMARY KEY (holdout_group_id, environment_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
