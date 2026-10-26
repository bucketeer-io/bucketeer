UPDATE holdout_group_feature
SET
    removed_at = $1
WHERE
    holdout_group_id = $2 AND environment_id = $3 AND feature_id = $4 AND removed_at = 0
