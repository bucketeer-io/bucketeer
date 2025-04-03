INSERT INTO ops_count (
	id,
	auto_ops_rule_id,
	clause_id,
	updated_at,
	ops_event_count,
	evaluation_count,
	feature_id,
	environment_id
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
	auto_ops_rule_id = VALUES(auto_ops_rule_id),
	clause_id = VALUES(clause_id),
	updated_at = VALUES(updated_at),
	ops_event_count = VALUES(ops_event_count),
	evaluation_count = VALUES(evaluation_count),
	feature_id = VALUES(feature_id)