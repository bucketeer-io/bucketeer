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
	$1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (id, environment_id) DO UPDATE SET
	auto_ops_rule_id = EXCLUDED.auto_ops_rule_id,
	clause_id = EXCLUDED.clause_id,
	updated_at = EXCLUDED.updated_at,
	ops_event_count = EXCLUDED.ops_event_count,
	evaluation_count = EXCLUDED.evaluation_count,
	feature_id = EXCLUDED.feature_id
