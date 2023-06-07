// package: bucketeer.eventcounter
// file: proto/eventcounter/goal_result.proto

import * as jspb from "google-protobuf";
import * as proto_eventcounter_variation_result_pb from "../../proto/eventcounter/variation_result_pb";

export class GoalResult extends jspb.Message {
  getGoalId(): string;
  setGoalId(value: string): void;

  clearVariationResultsList(): void;
  getVariationResultsList(): Array<proto_eventcounter_variation_result_pb.VariationResult>;
  setVariationResultsList(value: Array<proto_eventcounter_variation_result_pb.VariationResult>): void;
  addVariationResults(value?: proto_eventcounter_variation_result_pb.VariationResult, index?: number): proto_eventcounter_variation_result_pb.VariationResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalResult.AsObject;
  static toObject(includeInstance: boolean, msg: GoalResult): GoalResult.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GoalResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GoalResult;
  static deserializeBinaryFromReader(message: GoalResult, reader: jspb.BinaryReader): GoalResult;
}

export namespace GoalResult {
  export type AsObject = {
    goalId: string,
    variationResultsList: Array<proto_eventcounter_variation_result_pb.VariationResult.AsObject>,
  }
}

