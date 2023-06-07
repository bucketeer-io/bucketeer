// package: bucketeer.feature
// file: proto/feature/evaluation.proto

import * as jspb from "google-protobuf";
import * as proto_feature_variation_pb from "../../proto/feature/variation_pb";
import * as proto_feature_reason_pb from "../../proto/feature/reason_pb";

export class Evaluation extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getUserId(): string;
  setUserId(value: string): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  hasVariation(): boolean;
  clearVariation(): void;
  getVariation(): proto_feature_variation_pb.Variation | undefined;
  setVariation(value?: proto_feature_variation_pb.Variation): void;

  hasReason(): boolean;
  clearReason(): void;
  getReason(): proto_feature_reason_pb.Reason | undefined;
  setReason(value?: proto_feature_reason_pb.Reason): void;

  getVariationValue(): string;
  setVariationValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Evaluation.AsObject;
  static toObject(includeInstance: boolean, msg: Evaluation): Evaluation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Evaluation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Evaluation;
  static deserializeBinaryFromReader(message: Evaluation, reader: jspb.BinaryReader): Evaluation;
}

export namespace Evaluation {
  export type AsObject = {
    id: string,
    featureId: string,
    featureVersion: number,
    userId: string,
    variationId: string,
    variation?: proto_feature_variation_pb.Variation.AsObject,
    reason?: proto_feature_reason_pb.Reason.AsObject,
    variationValue: string,
  }
}

export class UserEvaluations extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  clearEvaluationsList(): void;
  getEvaluationsList(): Array<Evaluation>;
  setEvaluationsList(value: Array<Evaluation>): void;
  addEvaluations(value?: Evaluation, index?: number): Evaluation;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  clearArchivedFeatureIdsList(): void;
  getArchivedFeatureIdsList(): Array<string>;
  setArchivedFeatureIdsList(value: Array<string>): void;
  addArchivedFeatureIds(value: string, index?: number): string;

  getForceUpdate(): boolean;
  setForceUpdate(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserEvaluations.AsObject;
  static toObject(includeInstance: boolean, msg: UserEvaluations): UserEvaluations.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserEvaluations, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserEvaluations;
  static deserializeBinaryFromReader(message: UserEvaluations, reader: jspb.BinaryReader): UserEvaluations;
}

export namespace UserEvaluations {
  export type AsObject = {
    id: string,
    evaluationsList: Array<Evaluation.AsObject>,
    createdAt: number,
    archivedFeatureIdsList: Array<string>,
    forceUpdate: boolean,
  }

  export interface StateMap {
    QUEUED: 0;
    PARTIAL: 1;
    FULL: 2;
  }

  export const State: StateMap;
}

