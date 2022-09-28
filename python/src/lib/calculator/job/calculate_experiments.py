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
from typing import Dict, List

import grpc
from lib.calculator.domain import experiment as experiment_domain
from lib.calculator.domain import experiment_result as experiment_result_domain
from lib.calculator.job import metrics
from lib.calculator.storage import \
    mysql_experiment_result as mysql_experiment_result_storage
from proto.environment import service_pb2 as environment_service_pb2
from proto.eventcounter import experiment_result_pb2, goal_result_pb2
from proto.eventcounter import service_pb2 as ec_service_pb2
from proto.eventcounter import (timeseries_pb2, variation_count_pb2,
                                variation_result_pb2)
from proto.experiment import experiment_pb2
from proto.experiment import service_pb2 as experiment_service_pb2
from proto.experiment.command_pb2 import (FinishExperimentCommand,
                                          StartExperimentCommand)

_PAGE_SIZE = 500


class ExperimentCalculator:
    def __init__(
        self,
        environment_stub,
        experiment_stub,
        event_counter_stub,
        mysql_client,
        binomial,
        normal_inverse_gamma,
        logger,
        grpc_timeout=180,
    ):
        self._grpc_timeout = grpc_timeout
        self._environment_stub = environment_stub
        self._experiment_stub = experiment_stub
        self._event_counter_stub = event_counter_stub
        self._mysql_experiment_result_storage = (
            mysql_experiment_result_storage.MySQLExperimentResultStorage(mysql_client)
        )
        self._logger = logger
        self._binomial_model = binomial
        self._normal_inverse_gamma = normal_inverse_gamma
        self._now = datetime.now

    def run(self):
        environments = self._list_environments()
        metrics.target_items_gauge.labels(metrics.typeEnvironment, "").set(
            len(environments)
        )
        now = self._now()
        for env in environments or []:
            self._logger.info(
                "experiment calculator: start calculator over environment",
                {"environmentNamespace": env.namespace},
            )
            experiments = self._list_experiments(env.namespace)
            metrics.target_items_gauge.labels(
                metrics.typeExperiment, env.namespace
            ).set(len(experiments))
            for e in experiments:
                de = experiment_domain.Experiment(e)
                if not de.is_calculated(now):
                    continue
                experiment_result = self._create_experiment_result(env.namespace, de.pb)
                self._mysql_experiment_result_storage.upsert_multi(
                    env.namespace,
                    [experiment_result_domain.ExperimentResult(experiment_result)],
                )
                self._update_experiment_status(env.namespace, de, now)

    def _create_experiment_result(
        self,
        environment_namespace: str,
        experiment: experiment_pb2.Experiment,
    ) -> experiment_result_pb2.ExperimentResult:
        variation_ids = []
        for v in experiment.variations:
            variation_ids.append(v.id)
        end_ats = self._list_end_at(
            experiment.start_at, experiment.stop_at, int(self._now().timestamp())
        )
        expr_result = self._experiment_result(experiment.id)
        for goal_id in experiment.goal_ids:
            goal_result = expr_result.goal_results.add()
            goal_result.goal_id = goal_id
            variation_results = []
            for v in experiment.variations:
                variation_results.append(self._variation_result(v.id))
            for timestamp in end_ats:
                eval_vc = self._get_evaluation_count(
                    environment_namespace,
                    experiment.start_at,
                    timestamp,
                    experiment.feature_id,
                    experiment.feature_version,
                    variation_ids,
                )
                goal_vc = self._get_goal_count(
                    environment_namespace,
                    experiment.start_at,
                    timestamp,
                    goal_id,
                    experiment.feature_id,
                    experiment.feature_version,
                    variation_ids,
                )
                gr = self._create_goal_result(experiment, eval_vc, goal_vc)
                self._append_variation_results(
                    timestamp, variation_results, gr.variation_results
                )
            goal_result.variation_results.extend(variation_results)
        return expr_result

    def _append_variation_results(
        self,
        timestamp: int,
        dst_vrs: List[variation_result_pb2.VariationResult],
        src_vrs: List[variation_result_pb2.VariationResult],
    ):
        dst_vrs = sorted(dst_vrs, key=lambda vr: vr.variation_id)
        src_vrs = sorted(src_vrs, key=lambda vr: vr.variation_id)
        for dst_vr, src_vr in zip(dst_vrs, src_vrs):
            dst_vr.experiment_count.CopyFrom(src_vr.experiment_count)
            dst_vr.evaluation_count.CopyFrom(src_vr.evaluation_count)
            dst_vr.cvr_prob.CopyFrom(src_vr.cvr_prob)
            dst_vr.cvr_prob_best.CopyFrom(src_vr.cvr_prob_best)
            dst_vr.cvr_prob_beat_baseline.CopyFrom(src_vr.cvr_prob_beat_baseline)
            dst_vr.evaluation_user_count_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[float(src_vr.evaluation_count.user_count)],
                )
            )
            dst_vr.evaluation_event_count_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[float(src_vr.evaluation_count.event_count)],
                )
            )
            dst_vr.goal_user_count_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[float(src_vr.experiment_count.user_count)],
                )
            )
            dst_vr.goal_event_count_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[float(src_vr.experiment_count.event_count)],
                )
            )
            dst_vr.goal_value_sum_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp], values=[src_vr.experiment_count.value_sum]
                )
            )
            dst_vr.cvr_median_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp], values=[src_vr.cvr_prob.median]
                )
            )
            dst_vr.cvr_percentile025_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp], values=[src_vr.cvr_prob.percentile025]
                )
            )
            dst_vr.cvr_percentile975_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp], values=[src_vr.cvr_prob.percentile975]
                )
            )
            cvr = 0.0
            if src_vr.evaluation_count.user_count != 0:
                cvr = float(
                    src_vr.experiment_count.user_count
                    / src_vr.evaluation_count.user_count
                )
            dst_vr.cvr_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(timestamps=[timestamp], values=[cvr])
            )
            value_per_user = 0.0
            if src_vr.experiment_count.user_count != 0:
                value_per_user = float(
                    src_vr.experiment_count.value_sum
                    / src_vr.experiment_count.user_count
                )
            dst_vr.goal_value_sum_per_user_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp], values=[value_per_user]
                )
            )

            dst_vr.goal_value_sum_per_user_prob.CopyFrom(
                src_vr.goal_value_sum_per_user_prob
            )
            dst_vr.goal_value_sum_per_user_prob_best.CopyFrom(
                src_vr.goal_value_sum_per_user_prob_best
            )
            dst_vr.goal_value_sum_per_user_prob_beat_baseline.CopyFrom(
                src_vr.goal_value_sum_per_user_prob_beat_baseline
            )
            dst_vr.goal_value_sum_per_user_median_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[src_vr.goal_value_sum_per_user_prob.median],
                )
            )
            dst_vr.goal_value_sum_per_user_percentile025_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[src_vr.goal_value_sum_per_user_prob.percentile025],
                )
            )
            dst_vr.goal_value_sum_per_user_percentile975_timeseries.MergeFrom(
                timeseries_pb2.Timeseries(
                    timestamps=[timestamp],
                    values=[src_vr.goal_value_sum_per_user_prob.percentile975],
                )
            )

    def _list_environments(self):
        try:
            environments = []
            cursor = ""
            while True:
                resp = self._environment_stub.ListEnvironments(
                    environment_service_pb2.ListEnvironmentsRequest(
                        page_size=_PAGE_SIZE, cursor=cursor
                    ),
                    self._grpc_timeout,
                )
                environments.extend(resp.environments)
                environmentSize = len(resp.environments)
                if environmentSize < _PAGE_SIZE:
                    return environments
                cursor = resp.cursor
        # FIXME: Here it will be changed soon, maybe.
        # https://github.com/grpc/grpc/issues/9270#issuecomment-398796613
        except grpc.RpcError as rpc_error_call:
            self._logger.error(
                "experiment calculator: list environments failed",
                {"code": rpc_error_call.code(), "details": rpc_error_call.details()},
            )
            raise

    def _list_experiments(
        self,
        environment_namespace,
    ) -> List[experiment_pb2.Experiment]:
        experiments = []
        cursor = ""
        stopped_at = self._now() - timedelta(days=2)
        try:
            while True:
                req = experiment_service_pb2.ListExperimentsRequest(
                    environment_namespace=environment_namespace,
                    statuses=[
                        experiment_pb2.Experiment.WAITING,
                        experiment_pb2.Experiment.RUNNING,
                    ],
                    page_size=_PAGE_SIZE,
                    cursor=cursor,
                )
                # from is reserved word
                setattr(req, "from", int(stopped_at.timestamp()))
                resp = self._experiment_stub.ListExperiments(req, self._grpc_timeout)

                experiments.extend(resp.experiments)
                if len(resp.experiments) < _PAGE_SIZE:
                    return experiments
                cursor = resp.cursor
        except grpc.RpcError as rpc_error_call:
            self._logger.error(
                "experiment calculator: list experiments failed",
                {"code": rpc_error_call.code(), "details": rpc_error_call.details()},
            )
            raise

    def _list_end_at(self, start_at: int, end_at: int, now: int) -> List[int]:
        end_at = end_at if end_at < now else now
        timestamps = []
        day = 24 * 60 * 60
        for ts in range(start_at + day, end_at, day):
            timestamps.append(ts)
        timestamps.append(end_at)
        return timestamps

    def _experiment_result(self, experiment_id: str):
        return experiment_result_pb2.ExperimentResult(
            id=experiment_id,
            experiment_id=experiment_id,
            updated_at=int(self._now().timestamp()),
        )

    def _variation_count(self, variation_id: str):
        vc = variation_count_pb2.VariationCount()
        vc.variation_id = variation_id
        user_ts = timeseries_pb2.Timeseries()
        vc.user_timeseries_count.CopyFrom(user_ts)
        event_ts = timeseries_pb2.Timeseries()
        vc.event_timeseries_count.CopyFrom(event_ts)
        value_sum_ts = timeseries_pb2.Timeseries()
        vc.value_sum_timeseries_count.CopyFrom(value_sum_ts)
        return vc

    def _variation_result(self, variation_id: str):
        vr = variation_result_pb2.VariationResult()
        vr.variation_id = variation_id

        eval_user_ts = timeseries_pb2.Timeseries()
        vr.evaluation_user_count_timeseries.CopyFrom(eval_user_ts)
        eval_event_ts = timeseries_pb2.Timeseries()
        vr.evaluation_event_count_timeseries.CopyFrom(eval_event_ts)

        goal_user_ts = timeseries_pb2.Timeseries()
        vr.goal_user_count_timeseries.CopyFrom(goal_user_ts)
        goal_event_ts = timeseries_pb2.Timeseries()
        vr.goal_event_count_timeseries.CopyFrom(goal_event_ts)
        goal_value_sum_ts = timeseries_pb2.Timeseries()
        vr.goal_value_sum_timeseries.CopyFrom(goal_value_sum_ts)
        return vr

    def _get_evaluation_count(
        self,
        environment_namespace: str,
        start_at: int,
        end_at: int,
        feature_id: str,
        feature_version: int,
        variation_ids: List[str],
    ) -> Dict[str, variation_count_pb2.VariationCount]:
        try:
            resp = self._event_counter_stub.GetEvaluationCountV2(
                ec_service_pb2.GetEvaluationCountV2Request(
                    environment_namespace=environment_namespace,
                    start_at=start_at,
                    end_at=end_at,
                    feature_id=feature_id,
                    feature_version=feature_version,
                    variation_ids=variation_ids,
                ),
                self._grpc_timeout,
            )
            variation_counts = {}
            for vc in resp.count.realtime_counts:
                variation_counts[vc.variation_id] = vc
            return variation_counts
        except grpc.RpcError as rpc_error_call:
            self._logger.error(
                "experiment calculator: get evaluation count failed",
                {"code": rpc_error_call.code(), "details": rpc_error_call.details()},
            )
            raise

    def _get_goal_count(
        self,
        environment_namespace: str,
        start_at: int,
        end_at: int,
        goal_id: str,
        feature_id: str,
        feature_version: int,
        variation_ids: List[str],
    ) -> Dict[str, variation_count_pb2.VariationCount]:
        try:
            resp = self._event_counter_stub.GetGoalCountV2(
                ec_service_pb2.GetGoalCountV2Request(
                    environment_namespace=environment_namespace,
                    start_at=start_at,
                    end_at=end_at,
                    goal_id=goal_id,
                    feature_id=feature_id,
                    feature_version=feature_version,
                    variation_ids=variation_ids,
                ),
                self._grpc_timeout,
            )
            variation_counts = {}
            for vc in resp.goal_counts.realtime_counts:
                variation_counts[vc.variation_id] = vc
            return variation_counts
        except grpc.RpcError as rpc_error_call:
            self._logger.error(
                "experiment calculator: get goal count failed",
                {"code": rpc_error_call.code(), "details": rpc_error_call.details()},
            )
            raise

    def _create_goal_result(
        self,
        experiment: experiment_pb2.Experiment,
        evaluation_counts: Dict[str, variation_count_pb2.VariationCount],
        goal_counts: Dict[str, variation_count_pb2.VariationCount],
    ) -> List[goal_result_pb2.GoalResult]:
        goal_result = self._calc_goal_result(
            evaluation_counts, goal_counts, experiment.base_variation_id
        )
        goal_result.goal_id = experiment.goal_id
        return goal_result

    def _calc_goal_result(
        self,
        evaluation_variation_counts: Dict[str, variation_count_pb2.VariationCount],
        goal_variation_counts: Dict[str, variation_count_pb2.VariationCount],
        base_vid: str,
    ) -> goal_result_pb2.GoalResult:
        vids = []
        goal_uc = []
        eval_uc = []
        vrs = {}
        value_means = []
        value_vars = []
        # if not set, use 0.
        baseline_idx = 0
        gr = goal_result_pb2.GoalResult()
        for i, vid in enumerate(goal_variation_counts.keys()):
            goal_variation_count = goal_variation_counts[vid]
            vr = gr.variation_results.add()
            vr.variation_id = goal_variation_count.variation_id
            vr.experiment_count.CopyFrom(goal_variation_count)
            vid = goal_variation_count.variation_id
            if vid not in evaluation_variation_counts.keys():
                self._logger.error(
                    "experiment calculator: vid not found", {"variationId": vid}
                )
                return None

            evaluation_variation_count = evaluation_variation_counts[vid]
            vids.append(vid)
            goal_uc.append(goal_variation_count.user_count)
            eval_uc.append(evaluation_variation_count.user_count)
            vr.evaluation_count.CopyFrom(evaluation_variation_count)
            value_means.append(goal_variation_count.value_sum_per_user_mean)
            value_vars.append(goal_variation_count.value_sum_per_user_variance)

            vrs[vid] = vr
            if base_vid == vid:
                baseline_idx = i

        # Skip the calculation if evaluation count is less than goal count.
        for i in range(len(eval_uc)):
            if eval_uc[i] < goal_uc[i]:
                return gr

        cvr_result = self._binomial_model.run(vids, goal_uc, eval_uc, baseline_idx)
        for vid, variation_result in cvr_result.items():
            vrs[vid].cvr_prob.CopyFrom(variation_result.cvr_prob)
            vrs[vid].cvr_prob_best.CopyFrom(variation_result.cvr_prob_best)
            vrs[vid].cvr_prob_beat_baseline.CopyFrom(
                variation_result.cvr_prob_beat_baseline
            )

        # Skip the calculation if values are zero.
        for i in range(len(vids)):
            if goal_uc[i] == 0 or value_means[i] == 0.0 or value_vars[i] == 0.0:
                return gr
        value_result = self._normal_inverse_gamma.run(
            vids, value_means, value_vars, goal_uc, baseline_idx
        )
        for vid, variation_result in value_result.items():
            vrs[vid].goal_value_sum_per_user_prob.CopyFrom(
                variation_result.goal_value_sum_per_user_prob
            )
            vrs[vid].goal_value_sum_per_user_prob_best.CopyFrom(
                variation_result.goal_value_sum_per_user_prob_best
            )
            vrs[vid].goal_value_sum_per_user_prob_beat_baseline.CopyFrom(
                variation_result.goal_value_sum_per_user_prob_beat_baseline
            )
        return gr

    def _start_experiment(
        self,
        environment_namespace: str,
        id: str,
    ) -> None:
        try:
            self._experiment_stub.StartExperiment(
                experiment_service_pb2.StartExperimentRequest(
                    environment_namespace=environment_namespace,
                    id=id,
                    command=StartExperimentCommand(),
                ),
                self._grpc_timeout,
            )
            return
        except grpc.RpcError as rpc_error_call:
            self._logger.error(
                "experiment calculator: start experiment failed",
                {"code": rpc_error_call.code(), "details": rpc_error_call.details()},
            )
            raise

    def _finish_experiment(
        self,
        environment_namespace: str,
        id: str,
    ) -> None:
        try:
            self._experiment_stub.FinishExperiment(
                experiment_service_pb2.FinishExperimentRequest(
                    environment_namespace=environment_namespace,
                    id=id,
                    command=FinishExperimentCommand(),
                ),
                self._grpc_timeout,
            )
            return
        except grpc.RpcError as rpc_error_call:
            self._logger.error(
                "experiment calculator: finish experiment failed",
                {"code": rpc_error_call.code(), "details": rpc_error_call.details()},
            )
            raise

    def _update_experiment_status(
        self,
        environment_namespace: str,
        de: experiment_domain.Experiment,
        now: datetime,
    ) -> None:
        if de.is_updated_to_finish(now):
            self._finish_experiment(environment_namespace, de.pb.id)
            return
        if de.is_updated_to_running(now):
            self._start_experiment(environment_namespace, de.pb.id)
            return
