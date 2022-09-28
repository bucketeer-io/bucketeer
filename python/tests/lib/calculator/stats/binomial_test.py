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

from logging import getLogger

import pytest
from lib.calculator.stats.binomial import Binomial


def test_calc_result(mocker):
    logger = getLogger(__name__)
    binomial = Binomial(logger)
    vrs = binomial.run(["vid1", "vid2"], [38, 51], [101, 99], 0)
    if len(vrs) != 2:
        pytest.fail("incorrect variation result length: {}".format(len(vrs)))

    vr1 = vrs["vid1"]
    assert vr1.cvr_prob.mean >= 0.37
    assert vr1.cvr_prob.mean <= 0.38
    assert vr1.cvr_prob.sd >= 0.045
    assert vr1.cvr_prob.sd <= 0.05
    assert vr1.cvr_prob.rhat >= 0.9
    assert vr1.cvr_prob.rhat <= 1.1
    assert len(vr1.cvr_prob.histogram.hist) == 100
    assert len(vr1.cvr_prob.histogram.bins) == 101
    assert vr1.cvr_prob_best.mean >= 0.024
    assert vr1.cvr_prob_best.mean <= 0.026
    assert vr1.cvr_prob_best.sd >= 0.15
    assert vr1.cvr_prob_best.sd <= 0.16
    assert vr1.cvr_prob_best.rhat >= 0.9
    assert vr1.cvr_prob_best.rhat <= 1.1
    assert vr1.cvr_prob_beat_baseline.mean == 0.0
    assert vr1.cvr_prob_beat_baseline.sd == 0.0
    assert vr1.cvr_prob_beat_baseline.rhat == 0.0

    vr2 = vrs["vid2"]
    assert vr2.cvr_prob.mean >= 0.49
    assert vr2.cvr_prob.mean <= 0.52
    assert vr2.cvr_prob.sd >= 0.045
    assert vr2.cvr_prob.sd <= 0.05
    assert vr2.cvr_prob.rhat >= 0.9
    assert vr2.cvr_prob.rhat <= 1.1
    assert len(vr1.cvr_prob.histogram.hist) == 100
    assert len(vr1.cvr_prob.histogram.bins) == 101
    assert vr2.cvr_prob_best.mean >= 0.97
    assert vr2.cvr_prob_best.mean <= 0.98
    assert vr2.cvr_prob_best.sd >= 0.15
    assert vr2.cvr_prob_best.sd <= 0.16
    assert vr2.cvr_prob_best.rhat >= 0.9
    assert vr2.cvr_prob_best.rhat <= 1.1
    assert vr2.cvr_prob_beat_baseline.mean >= 0.97
    assert vr2.cvr_prob_beat_baseline.mean <= 0.98
    assert vr2.cvr_prob_beat_baseline.sd >= 0.15
    assert vr2.cvr_prob_beat_baseline.sd <= 0.16
    assert vr2.cvr_prob_beat_baseline.rhat >= 0.9
    assert vr2.cvr_prob_beat_baseline.rhat <= 1.1


if __name__ == "__main__":
    raise SystemExit(pytest.main([__file__]))
