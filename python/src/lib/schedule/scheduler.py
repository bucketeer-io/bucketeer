import asyncio
from datetime import datetime

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from lib.schedule import metrics
from lib.schedule.job import Job, Status


class Scheduler:
    def __init__(self, jobs: Job, logger):
        self._INTERVAL = 1
        self._logger = logger
        self.scheduler = AsyncIOScheduler()
        self._job_statuses = {}
        for j in jobs:
            self.scheduler.add_job(
                self._wrapf(j),
                "cron",
                year="*",
                week="*",
                month=j.month,
                day=j.day,
                day_of_week=j.day_of_week,
                hour=j.hour,
                minute=j.minute,
                second=j.second,
                max_instances=1,
            )
            self._job_statuses[j.name] = Status.SUCCESS

    def _wrapf(self, job: Job):
        def f():
            try:
                metrics.job_started_counter.labels(job.name).inc()
                start = datetime.now()
                job.func()
                end = datetime.now()
                self._job_statuses[job.name] = Status.SUCCESS
                self._logger.info(
                    "scheduler: job succeeded, jobName", {"jobName": job.name}
                )
                metrics.job_finished_counter.labels(
                    job.name, metrics.CODE_SUCCESS
                ).inc()
                metrics.job_duration_histogram.labels(job.name).observe(
                    (end - start).total_seconds()
                )
            except Exception as e:
                metrics.job_finished_counter.labels(job.name, metrics.CODE_FAIL).inc()
                self._logger.error(
                    "scheduler: job failed, jobName",
                    {"jobName": job.name, "error": str(e)},
                )
                self._job_statuses[job.name] = Status.FAIL
                return e

        return f

    def check(self) -> bool:
        for status in self._job_statuses.values():
            if status is Status.FAIL:
                return False
        return True

    async def run(self):
        self.scheduler.start()
        try:
            while True:
                await asyncio.sleep(self._INTERVAL)
        except asyncio.CancelledError:
            self._logger.info("scheduler: CancelledError")
            self.scheduler.shutdown(wait=True)
