SELECT
  COUNT(1)
FROM
  tag
JOIN
    environment_v2 env ON tag.environment_id = env.id
%s
