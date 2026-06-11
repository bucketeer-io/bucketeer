SELECT
    ex.id,
    ex.goal_id,
    ex.feature_id,
    ex.feature_version,
    ex.variations,
    ex.start_at,
    ex.stop_at,
    ex.stopped,
    ex.stopped_at,
    ex.created_at,
    ex.updated_at,
    ex.archived,
    ex.deleted,
    ex.goal_ids,
    ex.name,
    ex.description,
    ex.base_variation_id,
    ex.maintainer,
    ex.status,
    (
        SELECT
            jsonb_agg(jsonb_build_object('id', goal.id, 'name', goal.name))
        FROM
            goal
        WHERE
            ex.goal_ids::jsonb @> to_jsonb(goal.id) AND
            goal.environment_id = ex.environment_id
    ) AS goals
FROM
    experiment ex
WHERE
    id = $1 AND
    environment_id = $2
