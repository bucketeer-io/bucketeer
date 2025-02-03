SELECT
    goal.id,
    goal.name,
    goal.description,
    goal.connection_type,
    goal.archived,
    goal.deleted,
    goal.created_at,
    goal.updated_at,
    CASE
        WHEN (
            SELECT COUNT(1)
            FROM experiment ex1
            WHERE
                ex1.environment_id = ? AND
                ex1.goal_ids LIKE concat("%%", goal.id, "%%")
        ) > 0 THEN TRUE ELSE FALSE
    END AS is_in_use_status,
    (
        SELECT JSON_ARRAYAGG(
            JSON_OBJECT(
                'id', ex2.id,
                'name', ex2.experiment_name, -- alias in the subquery
                'feature_id', ex2.feature_id,
                'feature_name', ex2.feature_name, -- alias in the subquery
                'status', ex2.status
            )
        )
        FROM (
            SELECT DISTINCT
                ex2.id,
                ex2.name AS experiment_name, -- alias experiment name
                ex2.feature_id,
                ex2.status,
                ft.name AS feature_name -- alias feature name
            FROM experiment ex2
            JOIN feature ft ON ex2.feature_id = ft.id
            WHERE JSON_CONTAINS(ex2.goal_ids, CONCAT('"', goal.id, '"'), '$')
        ) AS ex2 -- alias the subquery to avoid outer reference conflict
    ) AS experiments
FROM
    goal
%s %s %s %s
