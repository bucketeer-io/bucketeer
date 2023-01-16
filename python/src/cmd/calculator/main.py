# Copyright 2022 The Bucketeer Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import asyncio
import logging
import platform

from environs import Env
from lib.calculator.job import calculate_experiments
from lib.calculator.stats import binomial, normal_inverse_gamma
from lib.environment.stub import stub as environment_stub
from lib.eventcounter.stub import stub as event_counter_stub
from lib.experiment.stub import stub as experiment_stub
from lib.health import health
from lib.log.logger import Logger
from lib.metrics import server as metrics_server
from lib.rpc import rpc
from lib.schedule import job, scheduler
from lib.signal import signal_handler as sh
from lib.storage.mysql import client as mysql_client


async def main():
    env = Env()
    env.read_env()
    mysql_user = env.str("BUCKETEER_CALCULATOR_MYSQL_USER")
    mysql_pass = env.str("BUCKETEER_CALCULATOR_MYSQL_PASS")
    mysql_host = env.str("BUCKETEER_CALCULATOR_MYSQL_HOST")
    mysql_port = env.int("BUCKETEER_CALCULATOR_MYSQL_PORT")
    mysql_db_name = env.str("BUCKETEER_CALCULATOR_MYSQL_DB_NAME")
    environment_service = env.str(
        "BUCKETEER_CALCULATOR_ENVIRONMENT_SERVICE", "localhost:9000"
    )
    experiment_service = env.str(
        "BUCKETEER_CALCULATOR_EXPERIMENT_SERVICE", "localhost:9000"
    )
    event_counter_service = env.str(
        "BUCKETEER_CALCULATOR_EVENT_COUNTER_SERVICE", "localhost:9000"
    )
    port = env.int("BUCKETEER_CALCULATOR_PORT", 9090)
    metrics_port = env.int("BUCKETEER_CALCULATOR_METRICS_PORT", 9002)
    log_level = env.log_level("BUCKETEER_CALCULATOR_LOG_LEVEL", logging.INFO)
    service_token_path = env.str("BUCKETEER_CALCULATOR_SERVICE_TOKEN")
    cert_path = env.str("BUCKETEER_CALCULATOR_CERT")
    key_path = env.str("BUCKETEER_CALCULATOR_KEY")
    job_cron_hour = env.str("BUCKETEER_CALCULATOR_JOB_CRON_HOUR")
    job_cron_minute = env.str("BUCKETEER_CALCULATOR_JOB_CRON_MINUTE")
    job_cron_second = env.str("BUCKETEER_CALCULATOR_JOB_CRON_SECOND")

    telepresence_root = env.str("TELEPRESENCE_ROOT", "")
    if telepresence_root:
        service_token_path = telepresence_root + service_token_path
        cert_path = telepresence_root + cert_path
        key_path = telepresence_root + key_path

    logger = Logger("calculator", log_level)
    _logger = logger.logger

    envrStub = environment_stub.create_stub(
        environment_service, cert_path, service_token_path
    )
    exprStub = experiment_stub.create_stub(
        experiment_service, cert_path, service_token_path
    )
    ecStub = event_counter_stub.create_stub(
        event_counter_service, cert_path, service_token_path
    )

    mc = mysql_client.Client(
        mysql_user, mysql_pass, mysql_host, mysql_port, mysql_db_name
    )

    calculator = calculate_experiments.ExperimentCalculator(
        envrStub,
        exprStub,
        ecStub,
        mc,
        binomial.Binomial(_logger),
        normal_inverse_gamma.NormalInverseGamma(_logger),
        _logger,
    )

    jobs = [
        job.Job(
            "calculate_experiments",
            calculator.run,
            hour=job_cron_hour,
            minute=job_cron_minute,
            second=job_cron_second,
        ),
    ]
    sch = scheduler.Scheduler(jobs, _logger)

    server = rpc.Server(port, cert_path, key_path, _logger)
    signal_handler = sh.SignalHandler(_logger)
    metrics = metrics_server.Server(metrics_port, _logger)

    checks = [sch.check]
    checker = health.Checker(checks, server, _logger)

    tasks = [
        server.run(),
        checker.run(),
        signal_handler.run(),
        sch.run(),
        logger.run(),
        metrics.run(),
    ]
    _logger.info("app starts running", {"pythonVersion": platform.python_version()})
    await asyncio.wait(tasks, return_when=asyncio.FIRST_COMPLETED)


if __name__ == "__main__":
    asyncio.run(main())
