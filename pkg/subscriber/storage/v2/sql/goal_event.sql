INSERT INTO goal_event (
    id, 
    environment_id, 
    timestamp, 
    goal_id, 
    value, 
    user_id, 
    user_data, 
    tag, 
    source_id,
    feature_id, 
    feature_version, 
    variation_id, 
    reason
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 