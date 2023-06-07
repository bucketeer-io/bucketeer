// package: bucketeer.eventcounter
// file: proto/eventcounter/experiment_result.proto

import * as jspb from "google-protobuf";
import * as proto_eventcounter_goal_result_pb from "../../proto/eventcounter/goal_result_pb";

export class ExperimentResult extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getExperimentId(): string;
  setExperimentId(value: string): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  clearGoalResultsList(): void;
  getGoalResultsList(): Array<proto_eventcounter_goal_result_pb.GoalResult>;
  setGoalResultsList(value: Array<proto_eventcounter_goal_result_pb.GoalResult>): void;
  addGoalResults(value?: proto_eventcounter_goal_result_pb.GoalResult, index?: number): proto_eventcounter_goal_result_pb.GoalResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentResult.AsObject;
  static toObject(includeInstance: boolean, msg: ExperimentResult): ExperimentResult.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExperimentResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentResult;
  static deserializeBinaryFromReader(message: ExperimentResult, reader: jspb.BinaryReader): ExperimentResult;
}

export namespace ExperimentResult {
  export type AsObject = {
    id: string,
    experimentId: string,
    updatedAt: number,
    goalResultsList: Array<proto_eventcounter_goal_result_pb.GoalResult.AsObject>,
  }
}

