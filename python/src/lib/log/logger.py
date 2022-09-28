import asyncio
import logging
import sys

from google.cloud.logging.handlers import ContainerEngineHandler
from lib.log.formatter import Formatter


class Logger:
    def __init__(self, log_name: str, log_level: str):
        self._interval = 1
        self.logger = self._setup_logging(log_name, log_level)

    def _setup_logging(self, log_name: str, log_level: str):
        logger = logging.getLogger(log_name)
        handler = ContainerEngineHandler(log_name, stream=sys.stdout)
        handler.setFormatter(Formatter())
        logger.addHandler(handler)
        logger.setLevel(log_level)
        logger.propagate = False
        return logger

    async def run(self):
        try:
            while True:
                await asyncio.sleep(self._interval)
        except asyncio.CancelledError:
            self.logger.info("logger: CancelledError")
            logging.shutdown()
            # stackdriver loging waits up 5 seconds to flush logs
            await asyncio.sleep(5)
