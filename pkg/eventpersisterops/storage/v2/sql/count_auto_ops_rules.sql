select count(*)
from auto_ops_rule
where triggered_at = 0
  and deleted = 0
