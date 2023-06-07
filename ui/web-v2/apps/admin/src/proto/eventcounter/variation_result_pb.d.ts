// package: bucketeer.eventcounter
// file: proto/eventcounter/variation_result.proto

import * as jspb from "google-protobuf";
import * as proto_eventcounter_variation_count_pb from "../../proto/eventcounter/variation_count_pb";
import * as proto_eventcounter_distribution_summary_pb from "../../proto/eventcounter/distribution_summary_pb";
import * as proto_eventcounter_timeseries_pb from "../../proto/eventcounter/timeseries_pb";

export class VariationResult extends jspb.Message {
  getVariationId(): string;
  setVariationId(value: string): void;

  hasExperimentCount(): boolean;
  clearExperimentCount(): void;
  getExperimentCount(): proto_eventcounter_variation_count_pb.VariationCount | undefined;
  setExperimentCount(value?: proto_eventcounter_variation_count_pb.VariationCount): void;

  hasEvaluationCount(): boolean;
  clearEvaluationCount(): void;
  getEvaluationCount(): proto_eventcounter_variation_count_pb.VariationCount | undefined;
  setEvaluationCount(value?: proto_eventcounter_variation_count_pb.VariationCount): void;

  hasCvrProbBest(): boolean;
  clearCvrProbBest(): void;
  getCvrProbBest(): proto_eventcounter_distribution_summary_pb.DistributionSummary | undefined;
  setCvrProbBest(value?: proto_eventcounter_distribution_summary_pb.DistributionSummary): void;

  hasCvrProbBeatBaseline(): boolean;
  clearCvrProbBeatBaseline(): void;
  getCvrProbBeatBaseline(): proto_eventcounter_distribution_summary_pb.DistributionSummary | undefined;
  setCvrProbBeatBaseline(value?: proto_eventcounter_distribution_summary_pb.DistributionSummary): void;

  hasCvrProb(): boolean;
  clearCvrProb(): void;
  getCvrProb(): proto_eventcounter_distribution_summary_pb.DistributionSummary | undefined;
  setCvrProb(value?: proto_eventcounter_distribution_summary_pb.DistributionSummary): void;

  hasEvaluationUserCountTimeseries(): boolean;
  clearEvaluationUserCountTimeseries(): void;
  getEvaluationUserCountTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setEvaluationUserCountTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasEvaluationEventCountTimeseries(): boolean;
  clearEvaluationEventCountTimeseries(): void;
  getEvaluationEventCountTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setEvaluationEventCountTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalUserCountTimeseries(): boolean;
  clearGoalUserCountTimeseries(): void;
  getGoalUserCountTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalUserCountTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalEventCountTimeseries(): boolean;
  clearGoalEventCountTimeseries(): void;
  getGoalEventCountTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalEventCountTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalValueSumTimeseries(): boolean;
  clearGoalValueSumTimeseries(): void;
  getGoalValueSumTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalValueSumTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasCvrMedianTimeseries(): boolean;
  clearCvrMedianTimeseries(): void;
  getCvrMedianTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setCvrMedianTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasCvrPercentile025Timeseries(): boolean;
  clearCvrPercentile025Timeseries(): void;
  getCvrPercentile025Timeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setCvrPercentile025Timeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasCvrPercentile975Timeseries(): boolean;
  clearCvrPercentile975Timeseries(): void;
  getCvrPercentile975Timeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setCvrPercentile975Timeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasCvrTimeseries(): boolean;
  clearCvrTimeseries(): void;
  getCvrTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setCvrTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalValueSumPerUserTimeseries(): boolean;
  clearGoalValueSumPerUserTimeseries(): void;
  getGoalValueSumPerUserTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalValueSumPerUserTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalValueSumPerUserProb(): boolean;
  clearGoalValueSumPerUserProb(): void;
  getGoalValueSumPerUserProb(): proto_eventcounter_distribution_summary_pb.DistributionSummary | undefined;
  setGoalValueSumPerUserProb(value?: proto_eventcounter_distribution_summary_pb.DistributionSummary): void;

  hasGoalValueSumPerUserProbBest(): boolean;
  clearGoalValueSumPerUserProbBest(): void;
  getGoalValueSumPerUserProbBest(): proto_eventcounter_distribution_summary_pb.DistributionSummary | undefined;
  setGoalValueSumPerUserProbBest(value?: proto_eventcounter_distribution_summary_pb.DistributionSummary): void;

  hasGoalValueSumPerUserProbBeatBaseline(): boolean;
  clearGoalValueSumPerUserProbBeatBaseline(): void;
  getGoalValueSumPerUserProbBeatBaseline(): proto_eventcounter_distribution_summary_pb.DistributionSummary | undefined;
  setGoalValueSumPerUserProbBeatBaseline(value?: proto_eventcounter_distribution_summary_pb.DistributionSummary): void;

  hasGoalValueSumPerUserMedianTimeseries(): boolean;
  clearGoalValueSumPerUserMedianTimeseries(): void;
  getGoalValueSumPerUserMedianTimeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalValueSumPerUserMedianTimeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalValueSumPerUserPercentile025Timeseries(): boolean;
  clearGoalValueSumPerUserPercentile025Timeseries(): void;
  getGoalValueSumPerUserPercentile025Timeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalValueSumPerUserPercentile025Timeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  hasGoalValueSumPerUserPercentile975Timeseries(): boolean;
  clearGoalValueSumPerUserPercentile975Timeseries(): void;
  getGoalValueSumPerUserPercentile975Timeseries(): proto_eventcounter_timeseries_pb.Timeseries | undefined;
  setGoalValueSumPerUserPercentile975Timeseries(value?: proto_eventcounter_timeseries_pb.Timeseries): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationResult.AsObject;
  static toObject(includeInstance: boolean, msg: VariationResult): VariationResult.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VariationResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VariationResult;
  static deserializeBinaryFromReader(message: VariationResult, reader: jspb.BinaryReader): VariationResult;
}

export namespace VariationResult {
  export type AsObject = {
    variationId: string,
    experimentCount?: proto_eventcounter_variation_count_pb.VariationCount.AsObject,
    evaluationCount?: proto_eventcounter_variation_count_pb.VariationCount.AsObject,
    cvrProbBest?: proto_eventcounter_distribution_summary_pb.DistributionSummary.AsObject,
    cvrProbBeatBaseline?: proto_eventcounter_distribution_summary_pb.DistributionSummary.AsObject,
    cvrProb?: proto_eventcounter_distribution_summary_pb.DistributionSummary.AsObject,
    evaluationUserCountTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    evaluationEventCountTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalUserCountTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalEventCountTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalValueSumTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    cvrMedianTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    cvrPercentile025Timeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    cvrPercentile975Timeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    cvrTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalValueSumPerUserTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalValueSumPerUserProb?: proto_eventcounter_distribution_summary_pb.DistributionSummary.AsObject,
    goalValueSumPerUserProbBest?: proto_eventcounter_distribution_summary_pb.DistributionSummary.AsObject,
    goalValueSumPerUserProbBeatBaseline?: proto_eventcounter_distribution_summary_pb.DistributionSummary.AsObject,
    goalValueSumPerUserMedianTimeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalValueSumPerUserPercentile025Timeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
    goalValueSumPerUserPercentile975Timeseries?: proto_eventcounter_timeseries_pb.Timeseries.AsObject,
  }
}

