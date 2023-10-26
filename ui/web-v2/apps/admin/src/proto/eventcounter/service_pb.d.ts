// package: bucketeer.eventcounter
// file: proto/eventcounter/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_eventcounter_experiment_result_pb from "../../proto/eventcounter/experiment_result_pb";
import * as proto_eventcounter_timeseries_pb from "../../proto/eventcounter/timeseries_pb";
import * as proto_eventcounter_variation_count_pb from "../../proto/eventcounter/variation_count_pb";

export class GetExperimentEvaluationCountRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getStartAt(): number;
  setStartAt(value: number): void;

  getEndAt(): number;
  setEndAt(value: number): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  clearVariationIdsList(): void;
  getVariationIdsList(): Array<string>;
  setVariationIdsList(value: Array<string>): void;
  addVariationIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentEvaluationCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentEvaluationCountRequest): GetExperimentEvaluationCountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentEvaluationCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentEvaluationCountRequest;
  static deserializeBinaryFromReader(message: GetExperimentEvaluationCountRequest, reader: jspb.BinaryReader): GetExperimentEvaluationCountRequest;
}

export namespace GetExperimentEvaluationCountRequest {
  export type AsObject = {
    environmentNamespace: string,
    startAt: number,
    endAt: number,
    featureId: string,
    featureVersion: number,
    variationIdsList: Array<string>,
  }
}

export class GetExperimentEvaluationCountResponse extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  clearVariationCountsList(): void;
  getVariationCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setVariationCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addVariationCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentEvaluationCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentEvaluationCountResponse): GetExperimentEvaluationCountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentEvaluationCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentEvaluationCountResponse;
  static deserializeBinaryFromReader(message: GetExperimentEvaluationCountResponse, reader: jspb.BinaryReader): GetExperimentEvaluationCountResponse;
}

export namespace GetExperimentEvaluationCountResponse {
  export type AsObject = {
    featureId: string,
    featureVersion: number,
    variationCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
  }
}

export class GetEvaluationTimeseriesCountRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getTimeRange(): GetEvaluationTimeseriesCountRequest.TimeRangeMap[keyof GetEvaluationTimeseriesCountRequest.TimeRangeMap];
  setTimeRange(value: GetEvaluationTimeseriesCountRequest.TimeRangeMap[keyof GetEvaluationTimeseriesCountRequest.TimeRangeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEvaluationTimeseriesCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetEvaluationTimeseriesCountRequest): GetEvaluationTimeseriesCountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEvaluationTimeseriesCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEvaluationTimeseriesCountRequest;
  static deserializeBinaryFromReader(message: GetEvaluationTimeseriesCountRequest, reader: jspb.BinaryReader): GetEvaluationTimeseriesCountRequest;
}

export namespace GetEvaluationTimeseriesCountRequest {
  export type AsObject = {
    environmentNamespace: string,
    featureId: string,
    timeRange: GetEvaluationTimeseriesCountRequest.TimeRangeMap[keyof GetEvaluationTimeseriesCountRequest.TimeRangeMap],
  }

  export interface TimeRangeMap {
    UNKNOWN: 0;
    TWENTY_FOUR_HOURS: 1;
    SEVEN_DAYS: 2;
    FOURTEEN_DAYS: 3;
    THIRTY_DAYS: 4;
  }

  export const TimeRange: TimeRangeMap;
}

export class GetEvaluationTimeseriesCountResponse extends jspb.Message {
  clearUserCountsList(): void;
  getUserCountsList(): Array<proto_eventcounter_timeseries_pb.VariationTimeseries>;
  setUserCountsList(value: Array<proto_eventcounter_timeseries_pb.VariationTimeseries>): void;
  addUserCounts(value?: proto_eventcounter_timeseries_pb.VariationTimeseries, index?: number): proto_eventcounter_timeseries_pb.VariationTimeseries;

  clearEventCountsList(): void;
  getEventCountsList(): Array<proto_eventcounter_timeseries_pb.VariationTimeseries>;
  setEventCountsList(value: Array<proto_eventcounter_timeseries_pb.VariationTimeseries>): void;
  addEventCounts(value?: proto_eventcounter_timeseries_pb.VariationTimeseries, index?: number): proto_eventcounter_timeseries_pb.VariationTimeseries;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEvaluationTimeseriesCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetEvaluationTimeseriesCountResponse): GetEvaluationTimeseriesCountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEvaluationTimeseriesCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEvaluationTimeseriesCountResponse;
  static deserializeBinaryFromReader(message: GetEvaluationTimeseriesCountResponse, reader: jspb.BinaryReader): GetEvaluationTimeseriesCountResponse;
}

export namespace GetEvaluationTimeseriesCountResponse {
  export type AsObject = {
    userCountsList: Array<proto_eventcounter_timeseries_pb.VariationTimeseries.AsObject>,
    eventCountsList: Array<proto_eventcounter_timeseries_pb.VariationTimeseries.AsObject>,
  }
}

export class GetExperimentResultRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getExperimentId(): string;
  setExperimentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentResultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentResultRequest): GetExperimentResultRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentResultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentResultRequest;
  static deserializeBinaryFromReader(message: GetExperimentResultRequest, reader: jspb.BinaryReader): GetExperimentResultRequest;
}

export namespace GetExperimentResultRequest {
  export type AsObject = {
    environmentNamespace: string,
    experimentId: string,
  }
}

export class GetExperimentResultResponse extends jspb.Message {
  hasExperimentResult(): boolean;
  clearExperimentResult(): void;
  getExperimentResult(): proto_eventcounter_experiment_result_pb.ExperimentResult | undefined;
  setExperimentResult(value?: proto_eventcounter_experiment_result_pb.ExperimentResult): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentResultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentResultResponse): GetExperimentResultResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentResultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentResultResponse;
  static deserializeBinaryFromReader(message: GetExperimentResultResponse, reader: jspb.BinaryReader): GetExperimentResultResponse;
}

export namespace GetExperimentResultResponse {
  export type AsObject = {
    experimentResult?: proto_eventcounter_experiment_result_pb.ExperimentResult.AsObject,
  }
}

export class ListExperimentResultsRequest extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  hasFeatureVersion(): boolean;
  clearFeatureVersion(): void;
  getFeatureVersion(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setFeatureVersion(value?: google_protobuf_wrappers_pb.Int32Value): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExperimentResultsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListExperimentResultsRequest): ListExperimentResultsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListExperimentResultsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExperimentResultsRequest;
  static deserializeBinaryFromReader(message: ListExperimentResultsRequest, reader: jspb.BinaryReader): ListExperimentResultsRequest;
}

export namespace ListExperimentResultsRequest {
  export type AsObject = {
    featureId: string,
    featureVersion?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    environmentNamespace: string,
  }
}

export class ListExperimentResultsResponse extends jspb.Message {
  getResultsMap(): jspb.Map<string, proto_eventcounter_experiment_result_pb.ExperimentResult>;
  clearResultsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExperimentResultsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListExperimentResultsResponse): ListExperimentResultsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListExperimentResultsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExperimentResultsResponse;
  static deserializeBinaryFromReader(message: ListExperimentResultsResponse, reader: jspb.BinaryReader): ListExperimentResultsResponse;
}

export namespace ListExperimentResultsResponse {
  export type AsObject = {
    resultsMap: Array<[string, proto_eventcounter_experiment_result_pb.ExperimentResult.AsObject]>,
  }
}

export class GetExperimentGoalCountRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getStartAt(): number;
  setStartAt(value: number): void;

  getEndAt(): number;
  setEndAt(value: number): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  clearVariationIdsList(): void;
  getVariationIdsList(): Array<string>;
  setVariationIdsList(value: Array<string>): void;
  addVariationIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentGoalCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentGoalCountRequest): GetExperimentGoalCountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentGoalCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentGoalCountRequest;
  static deserializeBinaryFromReader(message: GetExperimentGoalCountRequest, reader: jspb.BinaryReader): GetExperimentGoalCountRequest;
}

export namespace GetExperimentGoalCountRequest {
  export type AsObject = {
    environmentNamespace: string,
    startAt: number,
    endAt: number,
    goalId: string,
    featureId: string,
    featureVersion: number,
    variationIdsList: Array<string>,
  }
}

export class GetExperimentGoalCountResponse extends jspb.Message {
  getGoalId(): string;
  setGoalId(value: string): void;

  clearVariationCountsList(): void;
  getVariationCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setVariationCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addVariationCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentGoalCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentGoalCountResponse): GetExperimentGoalCountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentGoalCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentGoalCountResponse;
  static deserializeBinaryFromReader(message: GetExperimentGoalCountResponse, reader: jspb.BinaryReader): GetExperimentGoalCountResponse;
}

export namespace GetExperimentGoalCountResponse {
  export type AsObject = {
    goalId: string,
    variationCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
  }
}

export class GetOpsEvaluationUserCountRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOpsRuleId(): string;
  setOpsRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOpsEvaluationUserCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOpsEvaluationUserCountRequest): GetOpsEvaluationUserCountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOpsEvaluationUserCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOpsEvaluationUserCountRequest;
  static deserializeBinaryFromReader(message: GetOpsEvaluationUserCountRequest, reader: jspb.BinaryReader): GetOpsEvaluationUserCountRequest;
}

export namespace GetOpsEvaluationUserCountRequest {
  export type AsObject = {
    environmentNamespace: string,
    opsRuleId: string,
    clauseId: string,
    featureId: string,
    featureVersion: number,
    variationId: string,
  }
}

export class GetOpsEvaluationUserCountResponse extends jspb.Message {
  getOpsRuleId(): string;
  setOpsRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getCount(): number;
  setCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOpsEvaluationUserCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOpsEvaluationUserCountResponse): GetOpsEvaluationUserCountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOpsEvaluationUserCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOpsEvaluationUserCountResponse;
  static deserializeBinaryFromReader(message: GetOpsEvaluationUserCountResponse, reader: jspb.BinaryReader): GetOpsEvaluationUserCountResponse;
}

export namespace GetOpsEvaluationUserCountResponse {
  export type AsObject = {
    opsRuleId: string,
    clauseId: string,
    count: number,
  }
}

export class GetOpsGoalUserCountRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOpsRuleId(): string;
  setOpsRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOpsGoalUserCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOpsGoalUserCountRequest): GetOpsGoalUserCountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOpsGoalUserCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOpsGoalUserCountRequest;
  static deserializeBinaryFromReader(message: GetOpsGoalUserCountRequest, reader: jspb.BinaryReader): GetOpsGoalUserCountRequest;
}

export namespace GetOpsGoalUserCountRequest {
  export type AsObject = {
    environmentNamespace: string,
    opsRuleId: string,
    clauseId: string,
    featureId: string,
    featureVersion: number,
    variationId: string,
  }
}

export class GetOpsGoalUserCountResponse extends jspb.Message {
  getOpsRuleId(): string;
  setOpsRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getCount(): number;
  setCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOpsGoalUserCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOpsGoalUserCountResponse): GetOpsGoalUserCountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOpsGoalUserCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOpsGoalUserCountResponse;
  static deserializeBinaryFromReader(message: GetOpsGoalUserCountResponse, reader: jspb.BinaryReader): GetOpsGoalUserCountResponse;
}

export namespace GetOpsGoalUserCountResponse {
  export type AsObject = {
    opsRuleId: string,
    clauseId: string,
    count: number,
  }
}

export class GetMAUCountRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getYearMonth(): string;
  setYearMonth(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMAUCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMAUCountRequest): GetMAUCountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMAUCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMAUCountRequest;
  static deserializeBinaryFromReader(message: GetMAUCountRequest, reader: jspb.BinaryReader): GetMAUCountRequest;
}

export namespace GetMAUCountRequest {
  export type AsObject = {
    environmentNamespace: string,
    yearMonth: string,
  }
}

export class GetMAUCountResponse extends jspb.Message {
  getEventCount(): number;
  setEventCount(value: number): void;

  getUserCount(): number;
  setUserCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMAUCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMAUCountResponse): GetMAUCountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMAUCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMAUCountResponse;
  static deserializeBinaryFromReader(message: GetMAUCountResponse, reader: jspb.BinaryReader): GetMAUCountResponse;
}

export namespace GetMAUCountResponse {
  export type AsObject = {
    eventCount: number,
    userCount: number,
  }
}

export class SummarizeMAUCountsRequest extends jspb.Message {
  getYearMonth(): string;
  setYearMonth(value: string): void;

  getIsFinished(): boolean;
  setIsFinished(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SummarizeMAUCountsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SummarizeMAUCountsRequest): SummarizeMAUCountsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SummarizeMAUCountsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SummarizeMAUCountsRequest;
  static deserializeBinaryFromReader(message: SummarizeMAUCountsRequest, reader: jspb.BinaryReader): SummarizeMAUCountsRequest;
}

export namespace SummarizeMAUCountsRequest {
  export type AsObject = {
    yearMonth: string,
    isFinished: boolean,
  }
}

export class SummarizeMAUCountsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SummarizeMAUCountsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SummarizeMAUCountsResponse): SummarizeMAUCountsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SummarizeMAUCountsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SummarizeMAUCountsResponse;
  static deserializeBinaryFromReader(message: SummarizeMAUCountsResponse, reader: jspb.BinaryReader): SummarizeMAUCountsResponse;
}

export namespace SummarizeMAUCountsResponse {
  export type AsObject = {
  }
}

