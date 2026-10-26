INSERT INTO holdout_group_feature (
    id,
    holdout_group_id,
    feature_id,
    environment_id,
    holdout_variation_id,
    added_at,
    removed_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
