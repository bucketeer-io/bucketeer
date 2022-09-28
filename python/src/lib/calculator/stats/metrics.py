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

from prometheus_client import Histogram

binomial_compile_duration_histogram = Histogram(
    "binomial_compile_duration_seconds",
    "Duration of binomial model compilation in seconds.",
    namespace="bucketeer",
    subsystem="calculator",
    buckets=(30.0, 60.0, 90.0, 120.0, 180.0, 300.0, float("inf")),
)

binomial_sampling_duration_histogram = Histogram(
    "binomial_sampling_duration_seconds",
    "Duration of binomial model sampling in seconds.",
    namespace="bucketeer",
    subsystem="calculator",
    buckets=(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 4.0, 10.0, float("inf")),
)

binomial_run_duration_histogram = Histogram(
    "binomial_run_duration_seconds",
    "Duration of binomial model run in seconds.",
    namespace="bucketeer",
    subsystem="calculator",
    buckets=(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 4.0, 10.0, float("inf")),
)

normal_inverse_gamma_run_duration_histogram = Histogram(
    "normal_inverse_gamma_run_duration_seconds",
    "Duration of normal inverse gamma model run in seconds.",
    namespace="bucketeer",
    subsystem="calculator",
    buckets=(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 4.0, 10.0, float("inf")),
)
