UPDATE flag_trigger
SET feature_id        = $1,
    type              = $2,
    action            = $3,
    description       = $4,
    trigger_count     = $5,
    last_triggered_at = $6,
    token             = $7,
    disabled          = $8,
    created_at        = $9,
    updated_at        = $10
WHERE id = $11
  AND environment_id = $12
