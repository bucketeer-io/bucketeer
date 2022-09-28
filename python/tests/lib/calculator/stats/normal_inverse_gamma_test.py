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
from logging import getLogger

import numpy as np
import pytest
from lib.calculator.stats.normal_inverse_gamma import NormalInverseGamma


def test_run(mocker):
    logger = getLogger(__name__)
    nig = NormalInverseGamma(logger)
    # Pseudo sample statistic amount.
    v1_mean, v1_sd = 12, 10
    v2_mean, v2_sd = 15, 12
    sample_num = 20000
    # Generate peudo samples
    v1 = np.random.normal(loc=v1_mean, scale=v1_sd, size=sample_num)
    v2 = np.random.normal(loc=v2_mean, scale=v2_sd, size=sample_num)
    # Round up to zero.
    v1 = np.where(v1 < 0, 0, v1)
    v2 = np.where(v2 < 0, 0, v2)

    vids = ["vid1", "vid2"]
    means = [v1.mean(), v2.mean()]
    vars = [v1.var(), v2.var()]
    sizes = [v1.size, v2.size]
    baseline_idx = 0
    vrs = nig.run(vids, means, vars, sizes, baseline_idx)
    if len(vrs) != 2:
        pytest.fail("incorrect variation result length: {}".format(len(vrs)))

    vr1 = vrs["vid1"]
    assert vr1.goal_value_sum_per_user_prob.median > 12
    assert vr1.goal_value_sum_per_user_prob.median < 13.5
    assert vr1.goal_value_sum_per_user_prob.percentile025 > -2.0
    assert vr1.goal_value_sum_per_user_prob.percentile025 < -1.0
    assert vr1.goal_value_sum_per_user_prob.percentile975 > 26.5
    assert vr1.goal_value_sum_per_user_prob.percentile975 < 28.0
    assert vr1.goal_value_sum_per_user_prob_best.mean > 0.4
    assert vr1.goal_value_sum_per_user_prob_best.mean < 0.5
    assert vr1.goal_value_sum_per_user_prob_beat_baseline.mean == 0.0

    vr2 = vrs["vid2"]
    assert vr2.goal_value_sum_per_user_prob.median > 15
    assert vr2.goal_value_sum_per_user_prob.median < 17
    assert vr2.goal_value_sum_per_user_prob.percentile025 > -6.0
    assert vr2.goal_value_sum_per_user_prob.percentile025 < -4.0
    assert vr2.goal_value_sum_per_user_prob.percentile975 > 36.0
    assert vr2.goal_value_sum_per_user_prob.percentile975 < 37.5
    assert vr2.goal_value_sum_per_user_prob_best.mean > 0.4
    assert vr2.goal_value_sum_per_user_prob_best.mean < 0.6
    assert vr2.goal_value_sum_per_user_prob_beat_baseline.mean > 0.4
    assert vr2.goal_value_sum_per_user_prob_beat_baseline.mean < 0.6


def test_calc_beat_baseline(mocker):
    p = namedtuple("p", "input expected")
    patterns = [
        p(
            input=np.array([1, 2]),
            expected=np.array([[0, 1]]),
        ),
    ]
    for ptn in patterns:
        nig = NormalInverseGamma(getLogger(__name__))
        actual = nig._calc_beat_baseline(ptn.input, 0)
        assert ptn.expected[0, 0] == actual[0, 0]
        assert ptn.expected[0, 1] == actual[0, 1]


def test_calc_best(mocker):
    p = namedtuple("p", "input expected")
    patterns = [
        p(
            input=np.array([1, 2]),
            expected=np.array([[0, 1]]),
        ),
        p(
            input=np.array([1, 1]),
            expected=np.array([[1, 1]]),
        ),
    ]
    for ptn in patterns:
        nig = NormalInverseGamma(getLogger(__name__))
        actual = nig._calc_best(ptn.input)
        assert ptn.expected[0, 0] == actual[0, 0]
        assert ptn.expected[0, 1] == actual[0, 1]


if __name__ == "__main__":
    raise SystemExit(pytest.main([__file__]))
