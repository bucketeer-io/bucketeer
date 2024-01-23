DELETE FROM
    ops_progressive_rollout
WHERE
    id = ? AND
    environment_namespace = ?
