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

from collections import namedtuple
from datetime import datetime

import pytest
from proto.experiment import experiment_pb2
from lib.calculator.domain.experiment import Experiment


def test_experiment():
    actual = Experiment(experiment_pb2.Experiment())
    assert isinstance(actual, Experiment)
    assert isinstance(actual.pb, experiment_pb2.Experiment)


def test_is_calculated():
    t1 = _get_datetime("2018-06-24 08:15:27.243860")
    t2 = _get_datetime("2018-06-25 08:15:27.243860")

    p = namedtuple("p", "msg start_at now expected")
    patterns = [
        p(
            msg="true: start_at == now",
            start_at=t1,
            now=t1,
            expected=True,
        ),
        p(
            msg="true: start_at < now",
            start_at=t1,
            now=t2,
            expected=True,
        ),
        p(
            msg="true: start_at > now",
            start_at=t2,
            now=t1,
            expected=False,
        ),
    ]
    for ptn in patterns:
        e = experiment_pb2.Experiment()
        e.start_at = int(ptn.start_at.timestamp())
        de = Experiment(e)
        assert ptn.expected == de.is_calculated(ptn.now)


def test_is_updated_to_running():
    t1 = _get_datetime("2018-06-24 08:15:27.243860")
    t2 = _get_datetime("2018-06-25 08:15:27.243860")

    p = namedtuple("p", "msg start_at status now expected")
    patterns = [
        p(
            msg="true: start_at == now",
            start_at=t1,
            status=experiment_pb2.Experiment.WAITING,
            now=t1,
            expected=True,
        ),
        p(
            msg="true: start_at < now",
            start_at=t1,
            status=experiment_pb2.Experiment.WAITING,
            now=t2,
            expected=True,
        ),
        p(
            msg="true: start_at > now",
            start_at=t2,
            status=experiment_pb2.Experiment.WAITING,
            now=t1,
            expected=False,
        ),
        p(
            msg="true: RUNNING",
            start_at=t1,
            status=experiment_pb2.Experiment.RUNNING,
            now=t2,
            expected=False,
        ),
    ]
    for ptn in patterns:
        e = experiment_pb2.Experiment()
        e.start_at = int(ptn.start_at.timestamp())
        e.status = ptn.status
        de = Experiment(e)
        assert ptn.expected == de.is_updated_to_running(ptn.now)


def test_is_updated_to_stop():
    t1 = _get_datetime("2018-06-24 08:15:27.243860")
    t2 = _get_datetime("2018-06-25 08:15:27.243860")
    t3 = _get_datetime("2018-06-26 08:15:27.243860")
    t4 = _get_datetime("2018-06-27 08:15:27.243860")

    p = namedtuple("p", "msg stop_at status now expected")
    patterns = [
        p(
            msg="true: stop_at + 2 days < now",
            stop_at=t1,
            status=experiment_pb2.Experiment.RUNNING,
            now=t4,
            expected=True,
        ),
        p(
            msg="true: stop_at + 2 days == now",
            stop_at=t1,
            status=experiment_pb2.Experiment.RUNNING,
            now=t3,
            expected=False,
        ),
        p(
            msg="true: stop_at < now < stop_at + 2 days",
            stop_at=t1,
            status=experiment_pb2.Experiment.RUNNING,
            now=t3,
            expected=False,
        ),
        p(
            msg="true: stop_at > now",
            stop_at=t2,
            status=experiment_pb2.Experiment.RUNNING,
            now=t1,
            expected=False,
        ),
        p(
            msg="true: STOPPED",
            stop_at=t2,
            status=experiment_pb2.Experiment.STOPPED,
            now=t1,
            expected=False,
        ),
    ]
    for ptn in patterns:
        e = experiment_pb2.Experiment()
        e.stop_at = int(ptn.stop_at.timestamp())
        e.status = ptn.status
        de = Experiment(e)
        assert ptn.expected == de.is_updated_to_finish(ptn.now)


def _get_datetime(d: str) -> datetime:
    format = "%Y-%m-%d %H:%M:%S.%f"
    return datetime.strptime(d, format)


if __name__ == "__main__":
    raise SystemExit(pytest.main([__file__]))
