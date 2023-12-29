INSERT INTO flag_trigger(id,
                         feature_id,
                         environment_namespace,
                         type,
                         action,
                         description,
                         trigger_count,
                         last_triggered_at,
                         uuid,
                         token,
                         disabled,
                         created_at,
                         updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)