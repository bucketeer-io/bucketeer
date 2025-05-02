// package: bucketeer.eventcounter
// file: proto/eventcounter/goal_result.proto

import * as jspb from 'google-protobuf';
import * as proto_eventcounter_variation_result_pb from '../../proto/eventcounter/variation_result_pb';

export class GoalResult extends jspb.Message {
  getGoalId(): string;
  setGoalId(value: string): void;

  clearVariationResultsList(): void;
  getVariationResultsList(): Array<proto_eventcounter_variation_result_pb.VariationResult>;
  setVariationResultsList(
    value: Array<proto_eventcounter_variation_result_pb.VariationResult>
  ): void;
  addVariationResults(
    value?: proto_eventcounter_variation_result_pb.VariationResult,
    index?: number
  ): proto_eventcounter_variation_result_pb.VariationResult;

  hasSummary(): boolean;
  clearSummary(): void;
  getSummary(): Summary | undefined;
  setSummary(value?: Summary): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalResult.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GoalResult
  ): GoalResult.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GoalResult,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GoalResult;
  static deserializeBinaryFromReader(
    message: GoalResult,
    reader: jspb.BinaryReader
  ): GoalResult;
}

export namespace GoalResult {
  export type AsObject = {
    goalId: string;
    variationResultsList: Array<proto_eventcounter_variation_result_pb.VariationResult.AsObject>;
    summary?: Summary.AsObject;
  };
}

export class Summary extends jspb.Message {
  clearBestVariationsList(): void;
  getBestVariationsList(): Array<Summary.Variation>;
  setBestVariationsList(value: Array<Summary.Variation>): void;
  addBestVariations(
    value?: Summary.Variation,
    index?: number
  ): Summary.Variation;

  getTotalEvaluationUserCount(): number;
  setTotalEvaluationUserCount(value: number): void;

  getTotalGoalUserCount(): number;
  setTotalGoalUserCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Summary.AsObject;
  static toObject(includeInstance: boolean, msg: Summary): Summary.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Summary,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Summary;
  static deserializeBinaryFromReader(
    message: Summary,
    reader: jspb.BinaryReader
  ): Summary;
}

export namespace Summary {
  export type AsObject = {
    bestVariationsList: Array<Summary.Variation.AsObject>;
    totalEvaluationUserCount: number;
    totalGoalUserCount: number;
  };

  export class Variation extends jspb.Message {
    getId(): string;
    setId(value: string): void;

    getProbability(): number;
    setProbability(value: number): void;

    getIsbest(): boolean;
    setIsbest(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Variation.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: Variation
    ): Variation.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: Variation,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): Variation;
    static deserializeBinaryFromReader(
      message: Variation,
      reader: jspb.BinaryReader
    ): Variation;
  }

  export namespace Variation {
    export type AsObject = {
      id: string;
      probability: number;
      isbest: boolean;
    };
  }
}
