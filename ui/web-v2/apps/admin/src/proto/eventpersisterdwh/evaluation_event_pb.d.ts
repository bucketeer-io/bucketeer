// package: bucketeer.eventcounter
// file: proto/eventpersisterdwh/evaluation_event.proto

import * as jspb from "google-protobuf";

export class EvaluationEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getUserData(): string;
  setUserData(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  getReason(): string;
  setReason(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  getSourceId(): string;
  setSourceId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluationEvent.AsObject;
  static toObject(includeInstance: boolean, msg: EvaluationEvent): EvaluationEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EvaluationEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EvaluationEvent;
  static deserializeBinaryFromReader(message: EvaluationEvent, reader: jspb.BinaryReader): EvaluationEvent;
}

export namespace EvaluationEvent {
  export type AsObject = {
    id: string,
    featureId: string,
    featureVersion: number,
    userData: string,
    userId: string,
    variationId: string,
    reason: string,
    tag: string,
    sourceId: string,
    environmentNamespace: string,
    timestamp: number,
  }
}

