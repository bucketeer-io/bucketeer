from typing import Callable
import asyncio
from typing import List
from grpc_health.v1 import health
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

Checks = List[Callable[[], bool]]


class Checker:
    def __init__(self, checks: Checks, server, logger):
        self._interval = 1
        self._logger = logger
        self._checks = checks
        self._register_server(server)

    def _register_server(self, server):
        self._health_service = health.HealthServicer()
        health_pb2_grpc.add_HealthServicer_to_server(
            self._health_service, server.server
        )

    def _setServing(self):
        self._health_service.set(
            service="", status=health_pb2.HealthCheckResponse.SERVING
        )

    def _setNotServing(self):
        self._health_service.set(
            service="", status=health_pb2.HealthCheckResponse.NOT_SERVING
        )

    async def run(self):
        try:
            while True:
                for check in self._checks:
                    feedback = check()
                    if feedback is not True:
                        self._setNotServing()
                        break
                self._setServing()
                await asyncio.sleep(self._interval)
        except asyncio.CancelledError:
            self._logger.info("checker: CancelledError")
            self._health_service.enter_graceful_shutdown()
