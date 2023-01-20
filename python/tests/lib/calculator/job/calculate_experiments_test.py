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
from datetime import datetime, timedelta
from logging import getLogger
from typing import List

import pytest
from proto.eventcounter import (
    distribution_summary_pb2,
    evaluation_count_pb2,
    experiment_count_pb2,
    histogram_pb2,
)
from proto.eventcounter import service_pb2 as eventcounter_service_pb2
from proto.eventcounter import timeseries_pb2, variation_count_pb2, variation_result_pb2
from proto.experiment import experiment_pb2
from proto.experiment import service_pb2 as experiment_service_pb2
from proto.feature import variation_pb2
from pytest_mock import MockerFixture
from lib.calculator.domain import experiment as experiment_domain
from lib.calculator.job.calculate_experiments import ExperimentCalculator


def test_get_evaluation_count(mocker):
    ec = _experimentCalculator()
    insmock = mocker.Mock()
    vr0 = variation_count_pb2.VariationCount(variation_id="vid0", user_count=0)
    vr1 = variation_count_pb2.VariationCount(variation_id="vid1", user_count=1)
    resp = eventcounter_service_pb2.GetExperimentEvaluationCountResponse(variation_counts=[vr0, vr1])
    insmock.GetExperimentEvaluationCount.return_value = resp
    mocker.patch.object(ec, "_event_counter_stub", insmock)
    actual = ec._get_evaluation_count("", 0, 0, "", 0, [])
    assert actual["vid0"].user_count == 0
    assert actual["vid1"].user_count == 1


def test_get_goal_count(mocker):
    ec = _experimentCalculator()
    insmock = mocker.Mock()
    vr0 = variation_count_pb2.VariationCount(variation_id="vid0", user_count=0)
    vr1 = variation_count_pb2.VariationCount(variation_id="vid1", user_count=1)
    resp = eventcounter_service_pb2.GetExperimentGoalCountResponse(goal_id="gid0", variation_counts=[vr0, vr1])
    insmock.GetExperimentGoalCount.return_value = resp

    mocker.patch.object(ec, "_event_counter_stub", insmock)
    actual = ec._get_goal_count("", 0, 0, "gid0", "", 0, [])

    assert actual["vid0"].user_count == 0
    assert actual["vid1"].user_count == 1


def test_create_experiment_result(mocker):
    ec = _experimentCalculator()

    ec = _experimentCalculator()
    ec_insmock = mocker.Mock()
    eval_vr0 = variation_count_pb2.VariationCount(
        variation_id="vid0",
        user_count=5,
        event_count=10,
        value_sum=4,
    )
    eval_vr1 = variation_count_pb2.VariationCount(
        variation_id="vid1",
        user_count=4,
        event_count=12,
        value_sum=7,
    )
    resp = eventcounter_service_pb2.GetExperimentEvaluationCountResponse(variation_counts=[eval_vr0, eval_vr1])
    ec_insmock.GetExperimentEvaluationCount.return_value = resp

    goal_vr0 = variation_count_pb2.VariationCount(
        variation_id="vid0",
        user_count=2,
        event_count=4,
        value_sum=1.2,
        value_sum_per_user_mean=1.2,
        value_sum_per_user_variance=0.5,
    )
    goal_vr1 = variation_count_pb2.VariationCount(
        variation_id="vid1",
        user_count=1,
        event_count=2,
        value_sum=3.4,
        value_sum_per_user_mean=2.3,
        value_sum_per_user_variance=0.6,
    )
    resp = eventcounter_service_pb2.GetExperimentGoalCountResponse(goal_id="gid", variation_counts=[goal_vr0, goal_vr1])
    ec_insmock.GetExperimentGoalCount.return_value = resp

    mocker.patch.object(ec, "_event_counter_stub", ec_insmock)

    # CVR.
    cvr_vr = variation_result_pb2.VariationResult()
    cvr_prob = distribution_summary_pb2.DistributionSummary()
    cvr_prob.mean = 0.456
    cvr_prob.sd = 4.56
    cvr_prob.rhat = 45.6
    cvr_prob.median = 456.7
    cvr_prob.percentile025 = 4567.8
    cvr_prob.percentile975 = 45678.9
    cvr_prob_histogram = histogram_pb2.Histogram()
    cvr_prob_histogram.hist[:] = [1, 2]
    cvr_prob_histogram.bins[:] = [1, 2, 3]
    cvr_prob.histogram.CopyFrom(cvr_prob_histogram)
    cvr_vr.cvr_prob.CopyFrom(cvr_prob)

    cvr_prob_best = distribution_summary_pb2.DistributionSummary()
    cvr_prob_best.mean = 1.1
    cvr_prob_best.sd = 2.2
    cvr_prob_best.rhat = 3.3
    cvr_vr.cvr_prob_best.CopyFrom(cvr_prob_best)

    cvr_prob_beat_baseline = distribution_summary_pb2.DistributionSummary()
    cvr_prob_beat_baseline.mean = 1.11
    cvr_prob_beat_baseline.sd = 2.22
    cvr_prob_beat_baseline.rhat = 3.33
    cvr_vr.cvr_prob_beat_baseline.CopyFrom(cvr_prob_beat_baseline)

    binomial_insmock = mocker.Mock()
    binomial_insmock.run.return_value = {"vid0": cvr_vr, "vid1": cvr_vr}
    mocker.patch.object(ec, "_binomial_model", binomial_insmock)

    # Value sum per user.
    value_vr = variation_result_pb2.VariationResult()
    value_prob = distribution_summary_pb2.DistributionSummary()
    value_prob.median = 456.78
    value_prob.percentile025 = 4567.89
    value_prob.percentile975 = 45678.99
    value_vr.goal_value_sum_per_user_prob.CopyFrom(value_prob)

    value_prob_best = distribution_summary_pb2.DistributionSummary()
    value_prob_best.mean = 1.11
    value_vr.goal_value_sum_per_user_prob_best.CopyFrom(value_prob_best)

    value_prob_beat_baseline = distribution_summary_pb2.DistributionSummary()
    value_prob_beat_baseline.mean = 1.111
    value_vr.goal_value_sum_per_user_prob_beat_baseline.CopyFrom(
        value_prob_beat_baseline
    )

    nig_insmock = mocker.Mock()
    nig_insmock.run.return_value = {"vid0": value_vr, "vid1": value_vr}
    mocker.patch.object(ec, "_normal_inverse_gamma", nig_insmock)

    now = datetime.now()
    ec._now = lambda: now
    experiment = _create_experiment()
    experiment_result = ec._create_experiment_result("ns", experiment)
    goal_result = experiment_result.goal_results[0]
    assert goal_result.goal_id == "gid"
    for vr in goal_result.variation_results:
        if vr.variation_id == "vid0":
            vr0 = vr
            continue
        if vr.variation_id == "vid1":
            vr1 = vr
            continue
        pytest.fail("unknown variation id: {}".format(vr.variation_id))

    assert vr0.variation_id == "vid0"
    assert vr1.variation_id == "vid1"

    # vr0
    assert vr0.evaluation_count.variation_id == "vid0"
    assert vr0.evaluation_count.user_count == 5
    assert vr0.evaluation_count.event_count == 10
    assert vr0.experiment_count.variation_id == "vid0"
    assert vr0.experiment_count.user_count == 2
    assert vr0.experiment_count.event_count == 4
    assert vr0.experiment_count.value_sum == 1.2
    assert vr0.cvr_prob.mean == 0.456
    assert vr0.cvr_prob.sd == 4.56
    assert vr0.cvr_prob.rhat == 45.6
    assert vr0.cvr_prob.median == 456.7
    assert vr0.cvr_prob.percentile025 == 4567.8
    assert vr0.cvr_prob.percentile975 == 45678.9
    assert vr0.cvr_prob.histogram.hist == [1, 2]
    assert vr0.cvr_prob.histogram.bins == [1, 2, 3]
    assert vr0.cvr_prob_best.mean == 1.1
    assert vr0.cvr_prob_best.sd == 2.2
    assert vr0.cvr_prob_best.rhat == 3.3
    assert vr0.cvr_prob_beat_baseline.mean == 1.11
    assert vr0.cvr_prob_beat_baseline.sd == 2.22
    assert vr0.cvr_prob_beat_baseline.rhat == 3.33
    assert vr0.evaluation_user_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.evaluation_user_count_timeseries.values == [5.0, 5.0]
    assert vr0.evaluation_event_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.evaluation_event_count_timeseries.values == [10.0, 10.0]
    assert vr0.goal_user_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_user_count_timeseries.values == [2.0, 2.0]
    assert vr0.goal_event_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_event_count_timeseries.values == [4.0, 4.0]
    assert vr0.goal_value_sum_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_value_sum_timeseries.values == [1.2, 1.2]
    assert vr0.cvr_median_timeseries.timestamps == [1 * 24 * 60 * 60, 2 * 24 * 60 * 60]
    assert vr0.cvr_median_timeseries.values == [456.7, 456.7]
    assert vr0.cvr_percentile025_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.cvr_percentile025_timeseries.values == [4567.8, 4567.8]
    assert vr0.cvr_percentile975_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.cvr_percentile975_timeseries.values == [45678.9, 45678.9]
    assert vr0.cvr_timeseries.timestamps == [1 * 24 * 60 * 60, 2 * 24 * 60 * 60]
    assert vr0.cvr_timeseries.values == [0.4, 0.4]
    assert vr0.goal_value_sum_per_user_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_value_sum_per_user_timeseries.values == [0.6, 0.6]
    assert vr0.goal_value_sum_per_user_prob.median == 456.78
    assert vr0.goal_value_sum_per_user_prob.percentile025 == 4567.89
    assert vr0.goal_value_sum_per_user_prob.percentile975 == 45678.99
    assert vr0.goal_value_sum_per_user_prob_best.mean == 1.11
    assert vr0.goal_value_sum_per_user_prob_beat_baseline.mean == 1.111
    assert vr0.goal_value_sum_per_user_median_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_value_sum_per_user_median_timeseries.values == [456.78, 456.78]
    assert vr0.goal_value_sum_per_user_percentile025_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_value_sum_per_user_percentile025_timeseries.values == [
        4567.89,
        4567.89,
    ]
    assert vr0.goal_value_sum_per_user_percentile975_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr0.goal_value_sum_per_user_percentile975_timeseries.values == [
        45678.99,
        45678.99,
    ]

    # vr1
    assert vr1.evaluation_count.variation_id == "vid1"
    assert vr1.evaluation_count.user_count == 4
    assert vr1.evaluation_count.event_count == 12
    assert vr1.experiment_count.variation_id == "vid1"
    assert vr1.experiment_count.user_count == 1
    assert vr1.experiment_count.event_count == 2
    assert vr1.experiment_count.value_sum == 3.4
    assert vr1.cvr_prob.mean == 0.456
    assert vr1.cvr_prob.sd == 4.56
    assert vr1.cvr_prob.rhat == 45.6
    assert vr1.cvr_prob.median == 456.7
    assert vr1.cvr_prob.percentile025 == 4567.8
    assert vr1.cvr_prob.percentile975 == 45678.9
    assert vr1.cvr_prob.histogram.hist == [1, 2]
    assert vr1.cvr_prob.histogram.bins == [1, 2, 3]
    assert vr1.cvr_prob_best.mean == 1.1
    assert vr1.cvr_prob_best.sd == 2.2
    assert vr1.cvr_prob_best.rhat == 3.3
    assert vr1.cvr_prob_beat_baseline.mean == 1.11
    assert vr1.cvr_prob_beat_baseline.sd == 2.22
    assert vr1.cvr_prob_beat_baseline.rhat == 3.33
    assert vr1.evaluation_user_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.evaluation_user_count_timeseries.values == [4.0, 4.0]
    assert vr1.evaluation_event_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.evaluation_event_count_timeseries.values == [12.0, 12.0]
    assert vr1.goal_user_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_user_count_timeseries.values == [1.0, 1.0]
    assert vr1.goal_event_count_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_event_count_timeseries.values == [2.0, 2.0]
    assert vr1.goal_value_sum_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_value_sum_timeseries.values == [3.4, 3.4]
    assert vr1.cvr_median_timeseries.timestamps == [1 * 24 * 60 * 60, 2 * 24 * 60 * 60]
    assert vr1.cvr_median_timeseries.values == [456.7, 456.7]
    assert vr1.cvr_percentile025_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.cvr_percentile025_timeseries.values == [4567.8, 4567.8]
    assert vr1.cvr_percentile975_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.cvr_percentile975_timeseries.values == [45678.9, 45678.9]
    assert vr1.cvr_timeseries.timestamps == [1 * 24 * 60 * 60, 2 * 24 * 60 * 60]
    assert vr1.cvr_timeseries.values == [0.25, 0.25]
    assert vr1.goal_value_sum_per_user_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_value_sum_per_user_timeseries.values == [3.4, 3.4]
    assert vr1.goal_value_sum_per_user_prob.median == 456.78
    assert vr1.goal_value_sum_per_user_prob.percentile025 == 4567.89
    assert vr1.goal_value_sum_per_user_prob.percentile975 == 45678.99
    assert vr1.goal_value_sum_per_user_prob_best.mean == 1.11
    assert vr1.goal_value_sum_per_user_prob_beat_baseline.mean == 1.111
    assert vr1.goal_value_sum_per_user_median_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_value_sum_per_user_median_timeseries.values == [456.78, 456.78]
    assert vr1.goal_value_sum_per_user_percentile025_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_value_sum_per_user_percentile025_timeseries.values == [
        4567.89,
        4567.89,
    ]
    assert vr1.goal_value_sum_per_user_percentile975_timeseries.timestamps == [
        1 * 24 * 60 * 60,
        2 * 24 * 60 * 60,
    ]
    assert vr1.goal_value_sum_per_user_percentile975_timeseries.values == [
        45678.99,
        45678.99,
    ]


def test_update_experiment_status(mocker: MockerFixture):
    ec = _experimentCalculator()
    m = mocker.Mock()
    resp = experiment_service_pb2.StartExperimentResponse()
    m.StartExperiment.return_value = resp
    resp = experiment_service_pb2.FinishExperimentResponse()
    m.FinishExperiment.return_value = resp
    mocker.patch.object(ec, "_experiment_stub", m)

    now = datetime.now()
    now_unix = int(now.timestamp())
    stop_at_unix = int((now - timedelta(days=3)).timestamp())
    e = experiment_pb2.Experiment(
        status=experiment_pb2.Experiment.WAITING,
        stop_at=stop_at_unix,
    )
    de = experiment_domain.Experiment(e)
    ec._update_experiment_status("en", de, now)
    m.FinishExperiment.assert_called_once()

    now_unix = int(now.timestamp())
    stop_at_unix = int((now + timedelta(days=1)).timestamp())
    e = experiment_pb2.Experiment(
        status=experiment_pb2.Experiment.WAITING,
        start_at=now_unix,
        stop_at=stop_at_unix,
    )
    de = experiment_domain.Experiment(e)
    ec._update_experiment_status("en", de, now)
    m.StartExperiment.assert_called_once()


def test_append_variation_results(mocker):
    p = namedtuple("p", "msg timestamp dst_vrs src_vrs expected")
    patterns = [
        p(
            msg="true: running",
            timestamp=2,
            dst_vrs=_create_variation_results(
                variation_ids=["vid0", "vid1"],
                timestamps=[1],
                eval_user_counts=[1, 10],
                eval_event_counts=[2, 20],
                goal_user_counts=[3, 30],
                goal_event_counts=[4, 40],
                goal_value_sums=[5.5, 50.5],
                cvr_prob_median=[6.6, 60.6],
                cvr_prob_percentile025=[7.7, 70.7],
                cvr_prob_percentile975=[8.8, 80.8],
                eval_user_tss=[[1], [10]],
                eval_event_tss=[[2], [20]],
                goal_user_tss=[[3], [30]],
                goal_event_tss=[[4], [40]],
                goal_value_sum_tss=[[5.5], [50.5]],
                cvr_medians_tss=[[6.6], [60.6]],
                cvr_percentile025_tss=[[7.7], [70.7]],
                cvr_percentile975_tss=[[8.8], [80.8]],
                goal_value_sum_per_user_prob_median=[2.3, 4.3],
                goal_value_sum_per_user_prob_percentile025=[0.22, 0.33],
                goal_value_sum_per_user_prob_percentile975=[0.44, 0.55],
                goal_value_sum_per_user_medians_tss=[[1.2], [2.3]],
                goal_value_sum_per_user_percentile025_tss=[[0.11], [0.45]],
                goal_value_sum_per_user_percentile975_tss=[[0.12], [0.56]],
            ),
            src_vrs=_create_variation_results(
                variation_ids=["vid0", "vid1"],
                timestamps=[2],
                eval_user_counts=[2, 20],
                eval_event_counts=[3, 30],
                goal_user_counts=[4, 40],
                goal_event_counts=[5, 50],
                goal_value_sums=[6.6, 60.6],
                cvr_prob_median=[7.7, 70.7],
                cvr_prob_percentile025=[8.8, 80.8],
                cvr_prob_percentile975=[9.9, 90.9],
                eval_user_tss=[[2], [20]],
                eval_event_tss=[[3], [30]],
                goal_user_tss=[[4], [40]],
                goal_event_tss=[[5], [50]],
                goal_value_sum_tss=[[6.6], [60.6]],
                cvr_medians_tss=[[7.7], [70.7]],
                cvr_percentile025_tss=[[8.8], [80.8]],
                cvr_percentile975_tss=[[9.9], [90.9]],
                goal_value_sum_per_user_prob_median=[4.3, 5.3],
                goal_value_sum_per_user_prob_percentile025=[0.22, 0.33],
                goal_value_sum_per_user_prob_percentile975=[0.44, 0.55],
                goal_value_sum_per_user_medians_tss=[[1.2], [2.3]],
                goal_value_sum_per_user_percentile025_tss=[[0.11], [0.45]],
                goal_value_sum_per_user_percentile975_tss=[[0.12], [0.56]],
            ),
            expected=_create_variation_results(
                variation_ids=["vid0", "vid1"],
                timestamps=[1, 2],
                eval_user_counts=[2, 20],
                eval_event_counts=[3, 30],
                goal_user_counts=[4, 40],
                goal_event_counts=[5, 50],
                goal_value_sums=[6.6, 60.6],
                cvr_prob_median=[7.7, 70.7],
                cvr_prob_percentile025=[8.8, 80.8],
                cvr_prob_percentile975=[9.9, 90.9],
                eval_user_tss=[[1, 2], [10, 20]],
                eval_event_tss=[[2, 3], [20, 30]],
                goal_user_tss=[[3, 4], [30, 40]],
                goal_event_tss=[[4, 5], [40, 50]],
                goal_value_sum_tss=[[5.5, 6.6], [50.5, 60.6]],
                cvr_medians_tss=[[6.6, 7.7], [60.6, 70.7]],
                cvr_percentile025_tss=[[7.7, 8.8], [70.7, 80.8]],
                cvr_percentile975_tss=[[8.8, 9.9], [80.8, 90.9]],
                goal_value_sum_per_user_prob_median=[4.3, 5.3],
                goal_value_sum_per_user_prob_percentile025=[0.22, 0.33],
                goal_value_sum_per_user_prob_percentile975=[0.44, 0.55],
                goal_value_sum_per_user_medians_tss=[[1.2, 4.3], [2.3], 5.3],
                goal_value_sum_per_user_percentile025_tss=[[0.11, 0.22], [0.45, 0.33]],
                goal_value_sum_per_user_percentile975_tss=[[0.12, 0.44], [0.56, 0.55]],
            ),
        ),
    ]
    for ptn in patterns:
        ec = _experimentCalculator()
        ec._append_variation_results(ptn.timestamp, ptn.dst_vrs, ptn.src_vrs)

        actual = sorted(ptn.dst_vrs, key=lambda vr: vr.variation_id)
        expected = sorted(ptn.expected, key=lambda vr: vr.variation_id)

        assert (
            expected[0].evaluation_count.user_count
            == actual[0].evaluation_count.user_count
        )
        assert (
            expected[0].evaluation_count.event_count
            == actual[0].evaluation_count.event_count
        )
        assert (
            expected[0].evaluation_count.value_sum
            == actual[0].evaluation_count.value_sum
        )

        assert (
            expected[0].experiment_count.user_count
            == actual[0].experiment_count.user_count
        )
        assert (
            expected[0].experiment_count.event_count
            == actual[0].experiment_count.event_count
        )
        assert (
            expected[0].experiment_count.value_sum
            == actual[0].experiment_count.value_sum
        )

        assert expected[0].cvr_prob.median == actual[0].cvr_prob.median
        assert expected[0].cvr_prob.percentile025 == actual[0].cvr_prob.percentile025
        assert expected[0].cvr_prob.percentile975 == actual[0].cvr_prob.percentile975

        assert (
            expected[0].evaluation_user_count_timeseries.timestamps
            == actual[0].evaluation_user_count_timeseries.timestamps
        )
        assert (
            expected[0].evaluation_user_count_timeseries.values
            == actual[0].evaluation_user_count_timeseries.values
        )

        assert (
            expected[0].evaluation_event_count_timeseries.timestamps
            == actual[0].evaluation_event_count_timeseries.timestamps
        )
        assert (
            expected[0].evaluation_event_count_timeseries.values
            == actual[0].evaluation_event_count_timeseries.values
        )

        assert (
            expected[0].goal_user_count_timeseries.timestamps
            == actual[0].goal_user_count_timeseries.timestamps
        )
        assert (
            expected[0].goal_user_count_timeseries.values
            == actual[0].goal_user_count_timeseries.values
        )

        assert (
            expected[0].goal_event_count_timeseries.timestamps
            == actual[0].goal_event_count_timeseries.timestamps
        )
        assert (
            expected[0].goal_event_count_timeseries.values
            == actual[0].goal_event_count_timeseries.values
        )

        assert (
            expected[0].goal_value_sum_timeseries.timestamps
            == actual[0].goal_value_sum_timeseries.timestamps
        )
        assert (
            expected[0].goal_value_sum_timeseries.values
            == actual[0].goal_value_sum_timeseries.values
        )

        assert (
            expected[0].cvr_median_timeseries.timestamps
            == actual[0].cvr_median_timeseries.timestamps
        )
        assert (
            expected[0].cvr_median_timeseries.values
            == actual[0].cvr_median_timeseries.values
        )

        assert (
            expected[0].cvr_percentile025_timeseries.timestamps
            == actual[0].cvr_percentile025_timeseries.timestamps
        )
        assert (
            expected[0].cvr_percentile025_timeseries.values
            == actual[0].cvr_percentile025_timeseries.values
        )

        assert (
            expected[0].cvr_percentile975_timeseries.timestamps
            == actual[0].cvr_percentile975_timeseries.timestamps
        )
        assert (
            expected[0].cvr_percentile975_timeseries.values
            == actual[0].cvr_percentile975_timeseries.values
        )

        assert (
            expected[0].cvr_timeseries.timestamps == actual[0].cvr_timeseries.timestamps
        )
        assert expected[0].cvr_timeseries.values == actual[0].cvr_timeseries.values

        assert (
            expected[0].goal_value_sum_per_user_timeseries.timestamps
            == actual[0].goal_value_sum_per_user_timeseries.timestamps
        )
        assert (
            expected[0].goal_value_sum_per_user_timeseries.values
            == actual[0].goal_value_sum_per_user_timeseries.values
        )

        assert (
            expected[0].goal_value_sum_per_user_prob.median
            == actual[0].goal_value_sum_per_user_prob.median
        )
        assert (
            expected[0].goal_value_sum_per_user_prob.percentile025
            == actual[0].goal_value_sum_per_user_prob.percentile025
        )
        assert (
            expected[0].goal_value_sum_per_user_prob.percentile975
            == actual[0].goal_value_sum_per_user_prob.percentile975
        )

        assert (
            expected[0].goal_value_sum_per_user_median_timeseries.timestamps
            == actual[0].goal_value_sum_per_user_median_timeseries.timestamps
        )
        assert (
            expected[0].goal_value_sum_per_user_median_timeseries.values
            == actual[0].goal_value_sum_per_user_median_timeseries.values
        )

        assert (
            expected[0].goal_value_sum_per_user_percentile025_timeseries.timestamps
            == actual[0].goal_value_sum_per_user_percentile025_timeseries.timestamps
        )
        assert (
            expected[0].goal_value_sum_per_user_percentile025_timeseries.values
            == actual[0].goal_value_sum_per_user_percentile025_timeseries.values
        )

        assert (
            expected[0].goal_value_sum_per_user_percentile975_timeseries.timestamps
            == actual[0].goal_value_sum_per_user_percentile975_timeseries.timestamps
        )
        assert (
            expected[0].goal_value_sum_per_user_percentile975_timeseries.values
            == actual[0].goal_value_sum_per_user_percentile975_timeseries.values
        )


def test_list_end_at(mocker):
    day = 24 * 60 * 60

    p = namedtuple("p", "msg start_at stop_at now expected")
    patterns = [
        p(
            msg="1 hour",
            start_at=0,
            stop_at=1 * 60 * 60,
            now=32508810000,
            expected=[3600],
        ),
        p(
            msg="23 hours",
            start_at=0,
            stop_at=23 * 60 * 60,
            now=32508810000,
            expected=[82800],
        ),
        p(
            msg="1 day",
            start_at=0,
            stop_at=24 * 60 * 60,
            now=32508810000,
            expected=[86400],
        ),
        p(
            msg="3 days",
            start_at=0,
            stop_at=300000,
            now=32508810000,
            expected=[day, 2 * day, 3 * day, 300000],
        ),
        p(
            msg="3 days 18 hours",
            start_at=1614848400,  # 2021-03-04 09:00:00Z
            stop_at=1615086000,  # 2021-03-07 03:00:00Z
            now=32508810000,
            expected=[1614934800, 1615021200, 1615086000],
        ),
        p(
            msg="now is earlier than end_at",
            start_at=1614848400,  # 2021-03-04 09:00:00Z
            stop_at=1615086000,  # 2021-03-07 03:00:00Z
            now=1614967200,  # 2021-03-06 03:00:00Z
            expected=[1614934800, 1614967200],
        ),
    ]
    for ptn in patterns:
        ec = _experimentCalculator()
        assert ptn.expected == ec._list_end_at(ptn.start_at, ptn.stop_at, ptn.now)


def _experimentCalculator():
    logger = getLogger(__name__)
    return ExperimentCalculator(None, None, None, None, None, None, logger)


def _create_experiment():
    expr = experiment_pb2.Experiment()
    expr.id = "eid"
    expr.start_at = 0
    expr.stop_at = 2 * 24 * 60 * 60
    expr.base_variation_id = "vid1"
    expr.goal_ids.extend(["gid"])
    expr.variations.extend(
        [
            variation_pb2.Variation(id="vid0"),
            variation_pb2.Variation(id="vid1"),
        ]
    )
    return expr


def _create_variation_count(
    variation_id: str,
    user_count: int = 0,
    event_count: int = 0,
    value_sum: int = 0,
) -> variation_count_pb2.VariationCount:
    vc = variation_count_pb2.VariationCount()
    vc.variation_id = variation_id
    vc.user_count = user_count
    vc.event_count = event_count
    vc.value_sum = value_sum
    return vc


def _create_variation_result(
    variation_id: str,
    timestamps: List[int],
    eval_user_count: int,
    eval_event_count: int,
    goal_user_count: int,
    goal_event_count: int,
    goal_value_sum: int,
    cvr_prob_median: float,
    cvr_prob_percentile025: float,
    cvr_prob_percentile975: float,
    eval_user_ts: List[int],
    eval_event_ts: List[int],
    goal_user_ts: List[int],
    goal_event_ts: List[int],
    goal_value_sum_ts: List[float],
    cvr_median_ts: List[float],
    cvr_percentile025_ts: List[float],
    cvr_percentile975_ts: List[float],
    goal_value_sum_per_user_prob_median: float,
    goal_value_sum_per_user_prob_percentile025: float,
    goal_value_sum_per_user_prob_percentile975: float,
    goal_value_sum_per_user_medians_ts: List[float],
    goal_value_sum_per_user_percentile025_ts: List[float],
    goal_value_sum_per_user_percentile975_ts: List[float],
) -> variation_count_pb2.VariationCount:
    vr = variation_result_pb2.VariationResult()
    vr.variation_id = variation_id

    eval_vc = _create_variation_count(variation_id, eval_user_count, eval_event_count)
    vr.evaluation_count.CopyFrom(eval_vc)
    goal_vc = _create_variation_count(
        variation_id, goal_user_count, goal_event_count, goal_value_sum
    )
    vr.experiment_count.CopyFrom(goal_vc)

    cvr_prob_ds = _create_distribution_summary(
        cvr_prob_median, cvr_prob_percentile025, cvr_prob_percentile975
    )
    vr.cvr_prob.CopyFrom(cvr_prob_ds)

    eval_u_ts = timeseries_pb2.Timeseries()
    eval_u_ts.timestamps.extend(timestamps)
    eval_u_ts.values.extend([float(i) for i in eval_user_ts])
    vr.evaluation_user_count_timeseries.CopyFrom(eval_u_ts)

    eval_e_ts = timeseries_pb2.Timeseries()
    eval_e_ts.timestamps.extend(timestamps)
    eval_e_ts.values.extend([float(i) for i in eval_event_ts])
    vr.evaluation_event_count_timeseries.CopyFrom(eval_e_ts)

    goal_u_ts = timeseries_pb2.Timeseries()
    goal_u_ts.timestamps.extend(timestamps)
    goal_u_ts.values.extend([float(i) for i in goal_user_ts])
    vr.goal_user_count_timeseries.CopyFrom(goal_u_ts)

    goal_e_ts = timeseries_pb2.Timeseries()
    goal_e_ts.timestamps.extend(timestamps)
    goal_e_ts.values.extend([float(i) for i in goal_event_ts])
    vr.goal_event_count_timeseries.CopyFrom(goal_e_ts)

    goal_v_ts = timeseries_pb2.Timeseries()
    goal_v_ts.timestamps.extend(timestamps)
    goal_v_ts.values.extend(goal_value_sum_ts)
    vr.goal_value_sum_timeseries.CopyFrom(goal_v_ts)

    cvr_m_ts = timeseries_pb2.Timeseries()
    cvr_m_ts.timestamps.extend(timestamps)
    cvr_m_ts.values.extend(cvr_median_ts)
    vr.cvr_median_timeseries.CopyFrom(cvr_m_ts)

    cvr_02_ts = timeseries_pb2.Timeseries()
    cvr_02_ts.timestamps.extend(timestamps)
    cvr_02_ts.values.extend(cvr_percentile025_ts)
    vr.cvr_percentile025_timeseries.CopyFrom(cvr_02_ts)

    cvr_97_ts = timeseries_pb2.Timeseries()
    cvr_97_ts.timestamps.extend(timestamps)
    cvr_97_ts.values.extend(cvr_percentile975_ts)
    vr.cvr_percentile975_timeseries.CopyFrom(cvr_97_ts)

    cvr_ts = timeseries_pb2.Timeseries()
    cvr_ts.timestamps.extend(timestamps)
    cvr_ts.values.extend(
        [
            goal_count / eval_count
            for goal_count, eval_count in zip(goal_user_ts, eval_user_ts)
        ]
    )
    vr.cvr_timeseries.CopyFrom(cvr_ts)

    value_sum_per_user_ts = timeseries_pb2.Timeseries()
    value_sum_per_user_ts.timestamps.extend(timestamps)
    value_sum_per_user_ts.values.extend(
        [
            value_sum / goal_count
            for value_sum, goal_count in zip(goal_value_sum_ts, goal_user_ts)
        ]
    )
    vr.goal_value_sum_per_user_timeseries.CopyFrom(value_sum_per_user_ts)

    value_sum_per_user_prob_ds = _create_distribution_summary(
        goal_value_sum_per_user_prob_median,
        goal_value_sum_per_user_prob_percentile025,
        goal_value_sum_per_user_prob_percentile975,
    )
    vr.goal_value_sum_per_user_prob.CopyFrom(value_sum_per_user_prob_ds)

    value_sum_per_user_median_ts = timeseries_pb2.Timeseries()
    value_sum_per_user_median_ts.timestamps.extend(timestamps)
    value_sum_per_user_median_ts.values.extend(goal_value_sum_per_user_medians_ts)
    vr.goal_value_sum_per_user_median_timeseries.CopyFrom(value_sum_per_user_median_ts)

    value_sum_per_user_02_ts = timeseries_pb2.Timeseries()
    value_sum_per_user_02_ts.timestamps.extend(timestamps)
    value_sum_per_user_02_ts.values.extend(goal_value_sum_per_user_percentile025_ts)
    vr.goal_value_sum_per_user_percentile025_timeseries.CopyFrom(
        value_sum_per_user_02_ts
    )

    value_sum_per_user_97_ts = timeseries_pb2.Timeseries()
    value_sum_per_user_97_ts.timestamps.extend(timestamps)
    value_sum_per_user_97_ts.values.extend(goal_value_sum_per_user_percentile975_ts)
    vr.goal_value_sum_per_user_percentile975_timeseries.CopyFrom(
        value_sum_per_user_97_ts
    )

    return vr


def _create_distribution_summary(
    median: float,
    percentile025: float,
    percentile975: float,
) -> distribution_summary_pb2.DistributionSummary:
    return distribution_summary_pb2.DistributionSummary(
        median=median,
        percentile025=percentile025,
        percentile975=percentile975,
    )


def _create_variation_results(
    variation_ids: List[str],
    timestamps: List[int],
    eval_user_counts: List[int],
    eval_event_counts: List[int],
    goal_user_counts: List[int],
    goal_event_counts: List[int],
    goal_value_sums: List[int],
    cvr_prob_median: List[float],
    cvr_prob_percentile025: List[float],
    cvr_prob_percentile975: List[float],
    eval_user_tss: List[List[int]],
    eval_event_tss: List[List[int]],
    goal_user_tss: List[List[int]],
    goal_event_tss: List[List[int]],
    goal_value_sum_tss: List[List[float]],
    cvr_medians_tss: List[List[float]],
    cvr_percentile025_tss: List[List[float]],
    cvr_percentile975_tss: List[List[float]],
    goal_value_sum_per_user_prob_median: List[float],
    goal_value_sum_per_user_prob_percentile025: List[float],
    goal_value_sum_per_user_prob_percentile975: List[float],
    goal_value_sum_per_user_medians_tss: List[List[float]],
    goal_value_sum_per_user_percentile025_tss: List[List[float]],
    goal_value_sum_per_user_percentile975_tss: List[List[float]],
) -> List[variation_result_pb2.VariationResult]:
    vrs = []
    for i in range(len(variation_ids)):
        vrs.append(
            _create_variation_result(
                variation_id=variation_ids[i],
                timestamps=timestamps,
                eval_user_count=eval_user_counts[i],
                eval_event_count=eval_event_counts[i],
                goal_user_count=goal_user_counts[i],
                goal_event_count=goal_event_counts[i],
                goal_value_sum=goal_value_sums[i],
                cvr_prob_median=cvr_prob_median[i],
                cvr_prob_percentile025=cvr_prob_percentile025[i],
                cvr_prob_percentile975=cvr_prob_percentile975[i],
                eval_user_ts=eval_user_tss[i],
                eval_event_ts=eval_event_tss[i],
                goal_user_ts=goal_user_tss[i],
                goal_event_ts=goal_event_tss[i],
                goal_value_sum_ts=goal_value_sum_tss[i],
                cvr_median_ts=cvr_medians_tss[i],
                cvr_percentile025_ts=cvr_percentile025_tss[i],
                cvr_percentile975_ts=cvr_percentile975_tss[i],
                goal_value_sum_per_user_prob_median=goal_value_sum_per_user_prob_median[
                    i
                ],
                goal_value_sum_per_user_prob_percentile025=goal_value_sum_per_user_prob_percentile025[
                    i
                ],
                goal_value_sum_per_user_prob_percentile975=goal_value_sum_per_user_prob_percentile975[
                    i
                ],
                goal_value_sum_per_user_medians_ts=goal_value_sum_per_user_medians_tss[
                    i
                ],
                goal_value_sum_per_user_percentile025_ts=goal_value_sum_per_user_percentile025_tss[
                    i
                ],
                goal_value_sum_per_user_percentile975_ts=goal_value_sum_per_user_percentile975_tss[
                    i
                ],
            )
        )
    return vrs


if __name__ == "__main__":
    raise SystemExit(pytest.main([__file__]))
