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

import json
from typing import List

from google.protobuf import json_format
from lib.calculator.domain import experiment_result as domain
from lib.storage.mysql import client as mysql_client


class MySQLExperimentResultStorage:
    def __init__(self, client: mysql_client.Client):
        self._client = client

    def upsert_multi(
        self, environment_namespace: str, domains: List[domain.ExperimentResult]
    ):
        conn = self._client.get_conn()
        with conn:
            with conn.cursor() as cursor:
                for d in domains:
                    sql = (
                        "INSERT INTO experiment_result "
                        "(id, experiment_id, updated_at, data, environment_namespace) "
                        "VALUES (%s, %s, %s, %s, %s) "
                        "ON DUPLICATE KEY UPDATE "
                        "experiment_id = VALUES(experiment_id), "
                        "updated_at = VALUES(updated_at), "
                        "data = VALUES(data)"
                    )
                    dic = json_format.MessageToDict(
                        message=d.pb, preserving_proto_field_name=True
                    )
                    data = json.dumps(dic, ensure_ascii=False)
                    cursor.execute(
                        sql,
                        (
                            d.pb.id,
                            d.pb.experiment_id,
                            d.pb.updated_at,
                            data,
                            environment_namespace,
                        ),
                    )
            conn.commit()
