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

from datetime import datetime, timedelta

from proto.experiment import experiment_pb2


class Experiment:
    def __init__(self, pb: experiment_pb2.Experiment):
        self.pb = pb

    def is_updated_to_running(self, now: datetime) -> bool:
        now_unix = int(now.timestamp())
        if self.pb.status == experiment_pb2.Experiment.WAITING:
            if self.pb.start_at <= now_unix:
                return True
        return False

    def is_updated_to_finish(self, now: datetime) -> bool:
        two_days_ago_unix = int((now - timedelta(days=2)).timestamp())
        if self.pb.status in [
            experiment_pb2.Experiment.WAITING,
            experiment_pb2.Experiment.RUNNING,
        ]:
            if self.pb.stop_at < two_days_ago_unix:
                return True
        return False

    def is_calculated(self, now: datetime) -> bool:
        now_unix = int(now.timestamp())
        if self.pb.start_at <= now_unix:
            return True
        return False
