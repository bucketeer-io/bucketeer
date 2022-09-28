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

import math
from collections import namedtuple
from datetime import datetime
from typing import Dict, List

import numpy as np
from lib.calculator.stats import metrics
from proto.eventcounter import distribution_summary_pb2, variation_result_pb2

_PRIOR_MEAN = 30
_PRIOR_VAR = 2
_PRIOR_SIZE = 20
_PRIOR_ALPHA = 10
_PRIOR_BETA = 1000

_Distr = namedtuple("Distr", ("mu", "nu", "alpha", "beta", "n"))


class NormalInverseGamma:
    def __init__(self, logger):
        self._logger = logger

    def run(
        self,
        vids: List[str],
        means: List[float],
        vars: List[float],
        sizes: List[int],
        baseline_idx: int,
        post_gen_num: int = 25000,
    ) -> Dict[str, variation_result_pb2.VariationResult]:
        start = datetime.now()
        variation_num = len(means)
        posteriors = []
        samples = np.zeros((variation_num, post_gen_num))
        for i in range(variation_num):
            post = self._calc_posterior(
                sizes[i],
                means[i],
                vars[i],
                _PRIOR_SIZE,
                _PRIOR_MEAN,
                _PRIOR_VAR,
                _PRIOR_ALPHA,
                _PRIOR_BETA,
            )
            posteriors.append(post)
            samples[i, :] = self._gen_rnormgamma(post_gen_num, post)
        best = np.zeros((variation_num, post_gen_num))
        beat_baseline = np.zeros((variation_num, post_gen_num))
        for i in range(post_gen_num):
            best[:, i] = self._calc_best(samples[:, i])
            beat_baseline[:, i] = self._calc_beat_baseline(samples[:, i], baseline_idx)
        prob_best = np.sum(best, axis=1) / post_gen_num
        prob_beat_baseline = np.sum(beat_baseline, axis=1) / post_gen_num

        variation_results = {}
        for i in range(variation_num):
            vr = variation_result_pb2.VariationResult()
            vr.goal_value_sum_per_user_prob.CopyFrom(
                self._create_value_sum_prob(samples[i])
            )
            vr.goal_value_sum_per_user_prob_best.CopyFrom(
                self._create_value_sum_prob_best(prob_best[i])
            )
            vr.goal_value_sum_per_user_prob_beat_baseline.CopyFrom(
                self._create_value_sum_prob_beat_baseline(prob_beat_baseline[i])
            )
            variation_results[vids[i]] = vr

        metrics.normal_inverse_gamma_run_duration_histogram.observe(
            (datetime.now() - start).total_seconds()
        )
        return variation_results

    def _create_value_sum_prob(self, samples: List[float]):
        distr = distribution_summary_pb2.DistributionSummary()
        distr.median = np.median(samples)
        distr.percentile025 = np.percentile(samples, 2.5)
        distr.percentile975 = np.percentile(samples, 97.5)
        return distr

    def _create_value_sum_prob_best(self, prob_best: float):
        distr = distribution_summary_pb2.DistributionSummary()
        distr.mean = prob_best
        return distr

    def _create_value_sum_prob_beat_baseline(self, prob_beat_baseline: float):
        distr = distribution_summary_pb2.DistributionSummary()
        distr.mean = prob_beat_baseline
        return distr

    def _calc_posterior(
        self,
        this_n: int,
        this_mu: float,
        this_sigma: float,
        prior_n: int,
        prior_mu: float,
        prior_nu: float,
        prior_alpha: float,
        prior_beta: float,
    ):
        ret_n = this_n + prior_n
        # Take the logarithm to avoid from resulting big difference.
        n2 = math.log(this_n, 1.1)
        post_mu = (prior_nu * prior_mu + n2 * this_mu) / (prior_nu + n2)
        post_nu = prior_nu + n2
        post_alpha = prior_alpha + (n2 / 2)
        post_beta = (
            prior_beta
            + (1 / 2) * (this_sigma**2) * n2
            + (n2 * prior_nu / (prior_nu * n2)) * ((this_mu - prior_mu) ** 2) / 2
        )
        post = _Distr(post_mu, post_nu, post_alpha, post_beta, ret_n)
        return post

    def _gen_rnormgamma(self, n, posterior):
        return self._gen(n, posterior.mu, posterior.nu, posterior.alpha, posterior.beta)

    def _gen(self, n, mu, lmbd, alpha, beta):
        tau = 1 / np.random.gamma(alpha, scale=1 / beta, size=n)
        x = np.random.normal(loc=mu, scale=np.sqrt(tau / lmbd), size=n)
        return x

    def _calc_best(self, samples: List[float]):
        max = np.array([samples == samples.max()])
        return max.astype(int)

    def _calc_beat_baseline(self, samples: List[float], baseline_idx: int):
        baseline = samples[baseline_idx]
        beat_baseline = np.array([samples > baseline])
        return beat_baseline.astype(int)
