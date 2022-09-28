import asyncio

from prometheus_async import aio


class Server:
    def __init__(self, port, logger):
        self._port = port
        self._logger = logger

    async def run(self):
        self._logger.info("metrics: starting server on port", {"port": self._port})
        await aio.web.start_http_server(port=self._port)
        try:
            while True:
                await asyncio.sleep(1)
        except asyncio.CancelledError:
            self._logger.info("metrics: CancelledError")
