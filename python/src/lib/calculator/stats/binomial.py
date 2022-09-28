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
from datetime import datetime

import numpy as np
import pandas as pd
import pystan
from lib.calculator.stats import metrics
from proto.eventcounter import (distribution_summary_pb2, histogram_pb2,
                                variation_result_pb2)


class Binomial:
    def __init__(self, logger):
        self._logger = logger
        self._model = self._compile_model()

    def _compile_model(self):
        """
        Return the gcc compilation reslt.
        This process takes some time depending on machine resources.
        """
        model_code = """
        data {
            int<lower=0> g;
            int<lower=0> x[g];
            int<lower=0> n[g];
        }

        parameters {
            real<lower=0, upper=1> p[g];
        }

        model {
            for(i in 1:g){
                x[i] ~ binomial(n[i], p[i]);
            }
        }

        generated quantities {
            matrix[g, g] prob_upper;
            real prob_best[g];

            for(i in 1:g){
                real others[g-1];
                others = append_array(p[:i-1], p[i+1:]);
                prob_best[i] = p[i] > max(others) ? 1 : 0;
                for(j in 1:g){
                    prob_upper[i, j] = p[i] > p[j] ? 1 : 0;
                }
            }
        }
        """
        start = datetime.now()
        model = pystan.StanModel(model_code=model_code)
        end = datetime.now()
        metrics.binomial_compile_duration_histogram.observe(
            (end - start).total_seconds()
        )
        return model

    def run(
        self,
        vids,
        x,
        n,
        baseline_idx,
    ):
        # The index starts from 1 in PyStan.
        baseline_idx += 1
        start = datetime.now()
        num_variation = len(n)
        stan_data = {
            "g": num_variation,
            "x": np.array(x),
            "n": np.array(n),
        }
        par = [
            "p",
            "prob_upper",
            "prob_best",
        ]
        fit = self._model.sampling(
            data=stan_data,
            iter=21000,
            chains=5,
            warmup=1000,
            seed=1234,
            algorithm="NUTS",
        )
        metrics.binomial_sampling_duration_histogram.observe(
            (datetime.now() - start).total_seconds()
        )
        result = self._conv_fit(fit, par, vids, baseline_idx)
        metrics.binomial_run_duration_histogram.observe(
            (datetime.now() - start).total_seconds()
        )
        return result

    def _conv_fit(self, fit, par, vids, baseline_idx):
        variation_results = {}
        p_posterior_dist = np.array(fit.extract(permuted=True)["p"])
        summary = self._get_summary_df(fit, par)
        for i in range(1, len(vids) + 1):
            vr = variation_result_pb2.VariationResult()
            vr.cvr_prob.CopyFrom(self._create_cvr_prob(summary, p_posterior_dist, i))
            vr.cvr_prob_best.CopyFrom(self._create_cvr_prob_best(summary, i))
            vr.cvr_prob_beat_baseline.CopyFrom(
                self._create_cvr_prob_beat_baseline(summary, baseline_idx, i)
            )
            variation_results[vids[i - 1]] = vr
        return variation_results

    def _get_summary_df(self, fit, par):
        summary = fit.summary(pars=par)
        return pd.DataFrame(
            summary["summary"],
            index=summary["summary_rownames"],
            columns=summary["summary_colnames"],
        )

    def _create_cvr_prob(self, summary, p_posterior_dist, idx):
        prob = distribution_summary_pb2.DistributionSummary()
        prob.mean = summary.loc["p[{}]".format(idx), "mean"]
        prob.sd = summary.loc["p[{}]".format(idx), "sd"]
        prob.rhat = summary.loc["p[{}]".format(idx), "Rhat"]
        samples = p_posterior_dist.T[idx - 1]
        prob.median = np.median(samples)
        prob.percentile025 = np.percentile(samples, 2.5)
        prob.percentile975 = np.percentile(samples, 97.5)
        hist, bins = np.histogram(samples, bins=100)
        prob_histogram = histogram_pb2.Histogram()
        prob_histogram.hist[:] = hist
        prob_histogram.bins[:] = bins
        prob.histogram.CopyFrom(prob_histogram)
        return prob

    def _create_cvr_prob_best(self, summary, idx):
        prob_best = distribution_summary_pb2.DistributionSummary()
        prob_best.mean = summary.loc["prob_best[{}]".format(idx), "mean"]
        prob_best.sd = summary.loc["prob_best[{}]".format(idx), "sd"]
        prob_best.rhat = summary.loc["prob_best[{}]".format(idx), "Rhat"]
        return prob_best

    def _create_cvr_prob_beat_baseline(self, summary, baseline_idx, idx):
        prob_beat_baseline = distribution_summary_pb2.DistributionSummary()
        if idx is baseline_idx:
            prob_beat_baseline.mean = 0.0
            prob_beat_baseline.sd = 0.0
            prob_beat_baseline.rhat = 0.0
        else:
            prob_beat_baseline.mean = summary.loc[
                "prob_upper[{},{}]".format(idx, baseline_idx), "mean"
            ]
            prob_beat_baseline.sd = summary.loc[
                "prob_upper[{},{}]".format(idx, baseline_idx), "sd"
            ]
            prob_beat_baseline.rhat = summary.loc[
                "prob_upper[{},{}]".format(idx, baseline_idx), "Rhat"
            ]
        return prob_beat_baseline
