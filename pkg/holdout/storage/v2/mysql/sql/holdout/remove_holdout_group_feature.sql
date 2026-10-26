UPDATE holdout_group_feature
SET
    removed_at = ?
WHERE
    holdout_group_id = ? AND environment_id = ? AND feature_id = ? AND removed_at = 0
