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
            JSON_ARRAYAGG(JSON_OBJECT('id', goal.id, 'name', goal.name))
        FROM
            goal
        WHERE json_contains(ex.goal_ids, concat('"', goal.id, '"'), '$')
    ) as goals
FROM
    experiment ex
WHERE
    id = ? AND
    environment_id = ?