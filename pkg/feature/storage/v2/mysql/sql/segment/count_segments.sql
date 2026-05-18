SELECT
  COUNT(
    CASE
      WHEN (
        SELECT
          COUNT(1)
        FROM
          feature as ft
        WHERE
          ft.environment_id = seg.environment_id AND
          ft.rules LIKE concat("%%", seg.id, "%%")
      ) %s
    END
  )
FROM
  segment as seg
%s
