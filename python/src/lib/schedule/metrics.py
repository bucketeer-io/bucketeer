from prometheus_client import Histogram
from prometheus_client import Counter


CODE_SUCCESS = "Success"
CODE_FAIL = "Fail"

job_duration_histogram = Histogram(
    "scheduler_job_duration_seconds",
    "Job Duration.",
    labelnames=["name"],
    namespace="bucketeer",
    subsystem="calculator",
    buckets=(0.1, 1.0, 5.0, 10.0, 20.0, 40.0, 60.0, 120.0, float("inf")),
)

job_started_counter = Counter(
    "scheduler_started_jobs_total",
    "Total number of started jobs.",
    labelnames=["name"],
    namespace="bucketeer",
    subsystem="calculator",
)

job_finished_counter = Counter(
    "scheduler_finished_jobs_total",
    "Total number of finished jobs.",
    labelnames=["name", "code"],
    namespace="bucketeer",
    subsystem="calculator",
)
