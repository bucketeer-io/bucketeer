UPDATE flag_trigger
SET feature_id        = ?,
    type              = ?,
    action            = ?,
    description       = ?,
    trigger_count     = ?,
    last_triggered_at = ?,
    uuid              = ?,
    token             = ?,
    disabled          = ?,
    created_at        = ?,
    updated_at        = ?
WHERE id = ?
  AND environment_namespace = ?
