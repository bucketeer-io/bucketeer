import signal
import asyncio


class SignalHandler:
    def __init__(self, logger):
        self._INTERVAL = 1
        self._logger = logger
        self._kill_now = False
        self._set_handler(self._exit_gracefully)

    def _set_handler(self, handler):
        signal.signal(signal.SIGINT, handler)
        signal.signal(signal.SIGTERM, handler)

    def _exit_gracefully(self, signum, frame):
        self._logger.info("signal handler: Signal reveived")
        self._kill_now = True
        # No longer accept any signals
        # "*_" means that handler should accept 2 parameters.
        self._set_handler(lambda *_: None)

    async def check(self, interval: int):
        try:
            while True:
                if self._kill_now:
                    self._logger.info("signal handler: stop")
                    return
                await asyncio.sleep(interval)
        except asyncio.CancelledError:
            self._logger.debug("signal handler: CancelledError")

    async def run(self):
        await self.check(self._INTERVAL)
