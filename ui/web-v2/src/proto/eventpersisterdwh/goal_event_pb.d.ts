// package: bucketeer.eventcounter
// file: proto/eventpersisterdwh/goal_event.proto

import * as jspb from "google-protobuf";

export class GoalEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getValue(): number;
  setValue(value: number): void;

  getUserData(): string;
  setUserData(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  getSourceId(): string;
  setSourceId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  getReason(): string;
  setReason(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalEvent.AsObject;
  static toObject(includeInstance: boolean, msg: GoalEvent): GoalEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GoalEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GoalEvent;
  static deserializeBinaryFromReader(message: GoalEvent, reader: jspb.BinaryReader): GoalEvent;
}

export namespace GoalEvent {
  export type AsObject = {
    id: string,
    goalId: string,
    value: number,
    userData: string,
    userId: string,
    tag: string,
    sourceId: string,
    environmentNamespace: string,
    timestamp: number,
    featureId: string,
    featureVersion: number,
    variationId: string,
    reason: string,
  }
}

