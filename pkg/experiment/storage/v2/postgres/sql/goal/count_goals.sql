SELECT
    COUNT(
        CASE
            WHEN (
                SELECT COUNT(1)
                FROM experiment
                WHERE
                    experiment.environment_id = goal.environment_id AND
                    jsonb_exists(experiment.goal_ids, goal.id)
            ) %s
        END
    )
FROM
    goal
%s
