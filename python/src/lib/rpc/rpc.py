from concurrent import futures
import os
import grpc
import asyncio


class Server:
    def __init__(self, port, cert_path, key_path, logger):
        self._interval = 1
        self._port = port
        self._logger = logger
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        self._credentials = grpc.ssl_server_credentials(
            (
                (
                    self._load_credential_from_file(key_path),
                    self._load_credential_from_file(cert_path),
                ),
            )
        )

    @staticmethod
    def _load_credential_from_file(filepath):
        real_path = os.path.join(os.path.dirname(__file__), filepath)
        with open(real_path, "rb") as f:
            return f.read()

    async def run(self):
        self.server.add_secure_port("[::]:%d" % self._port, self._credentials)
        self.server.start()
        try:
            while True:
                await asyncio.sleep(self._interval)
        except asyncio.CancelledError:
            self._logger.info("server: CancelledError")
            self.server.stop(5)
