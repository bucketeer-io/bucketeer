INSERT INTO flag_triggers(id,
                          feature_id,
                          environment_namespace,
                          type,
                          action,
                          description,
                          trigger_times,
                          last_triggered_at,
                          uuid,
                          disabled,
                          created_at,
                          updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)