SELECT id,
       feature_id,
       type,
       `action`,
       description,
       trigger_count,
       last_triggered_at,
       disabled,
       created_at,
       updated_at
FROM flag_trigger %s %s %s