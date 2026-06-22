DELETE FROM scheduled_feature_change
WHERE id = $1
  AND environment_id = $2
