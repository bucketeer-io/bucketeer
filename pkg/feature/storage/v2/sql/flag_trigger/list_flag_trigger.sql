SELECT id,
       feature_id,
       environment_namespace,
       type,
       `action`,
       description,
       trigger_count,
       last_triggered_at,
       uuid,
       disabled,
       created_at,
       updated_at
FROM flag_trigger %s %s %s