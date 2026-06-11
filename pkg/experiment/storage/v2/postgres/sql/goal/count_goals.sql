SELECT
    COUNT(
        CASE
            WHEN (
                SELECT COUNT(1)
                FROM experiment
                WHERE
                    experiment.environment_id = goal.environment_id AND
                    experiment.goal_ids::jsonb @> to_jsonb(goal.id)
            ) %s
        END
    )
FROM
    goal
%s
