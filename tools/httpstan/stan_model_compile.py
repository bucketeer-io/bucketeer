#  Copyright 2024 The Bucketeer Authors.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
import asyncio
import typing

import aiohttp.web
import httpstan.app

with open("model.stan", "r") as f:
    model_code = f.read()

data = {
    "program_code": model_code,
}


async def server_start():
    app = httpstan.app.make_app()
    runner = aiohttp.web.AppRunner(app)
    await runner.setup()
    site = aiohttp.web.TCPSite(runner, "0.0.0.0", 8080)
    await site.start()


async def main():
    await server_start()
    async with aiohttp.ClientSession() as session:
        async with session.post("http://localhost:8080/v1/models", json=data) as resp:
            assert resp.status == 201
            response_payload = await resp.json()
    model_name = typing.cast(str, response_payload["name"])
    print(f"model_name: {model_name}")


if __name__ == '__main__':
    asyncio.run(main())
