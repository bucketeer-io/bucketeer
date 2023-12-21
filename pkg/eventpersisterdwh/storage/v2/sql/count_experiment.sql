select count(*)
from experiment
where deleted = 0
  and status = 1