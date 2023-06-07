// package: bucketeer.eventcounter
// file: proto/eventcounter/experiment_count.proto

import * as jspb from "google-protobuf";
import * as proto_eventcounter_variation_count_pb from "../../proto/eventcounter/variation_count_pb";

export class ExperimentCount extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  clearRealtimeCountsList(): void;
  getRealtimeCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setRealtimeCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addRealtimeCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  clearBatchCountsList(): void;
  getBatchCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setBatchCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addBatchCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  clearGoalCountsList(): void;
  getGoalCountsList(): Array<GoalCounts>;
  setGoalCountsList(value: Array<GoalCounts>): void;
  addGoalCounts(value?: GoalCounts, index?: number): GoalCounts;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentCount.AsObject;
  static toObject(includeInstance: boolean, msg: ExperimentCount): ExperimentCount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExperimentCount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentCount;
  static deserializeBinaryFromReader(message: ExperimentCount, reader: jspb.BinaryReader): ExperimentCount;
}

export namespace ExperimentCount {
  export type AsObject = {
    id: string,
    featureId: string,
    featureVersion: number,
    goalId: string,
    realtimeCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
    batchCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
    updatedAt: number,
    goalCountsList: Array<GoalCounts.AsObject>,
  }
}

export class GoalCounts extends jspb.Message {
  getGoalId(): string;
  setGoalId(value: string): void;

  clearRealtimeCountsList(): void;
  getRealtimeCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setRealtimeCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addRealtimeCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  clearBatchCountsList(): void;
  getBatchCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setBatchCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addBatchCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalCounts.AsObject;
  static toObject(includeInstance: boolean, msg: GoalCounts): GoalCounts.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GoalCounts, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GoalCounts;
  static deserializeBinaryFromReader(message: GoalCounts, reader: jspb.BinaryReader): GoalCounts;
}

export namespace GoalCounts {
  export type AsObject = {
    goalId: string,
    realtimeCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
    batchCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
  }
}

