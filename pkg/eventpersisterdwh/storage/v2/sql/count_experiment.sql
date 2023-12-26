select id, status, stop_at
from experiment
where deleted = 0
  and status in (1, 2)