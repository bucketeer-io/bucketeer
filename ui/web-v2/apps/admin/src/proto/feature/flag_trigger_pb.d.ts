// package: bucketeer.feature
// file: proto/feature/flag_trigger.proto

import * as jspb from "google-protobuf";

export class FlagTrigger extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getType(): FlagTrigger.TypeMap[keyof FlagTrigger.TypeMap];
  setType(value: FlagTrigger.TypeMap[keyof FlagTrigger.TypeMap]): void;

  getAction(): FlagTrigger.ActionMap[keyof FlagTrigger.ActionMap];
  setAction(value: FlagTrigger.ActionMap[keyof FlagTrigger.ActionMap]): void;

  getDescription(): string;
  setDescription(value: string): void;

  getTriggerCount(): number;
  setTriggerCount(value: number): void;

  getLastTriggeredAt(): number;
  setLastTriggeredAt(value: number): void;

  getUuid(): string;
  setUuid(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTrigger.AsObject;
  static toObject(includeInstance: boolean, msg: FlagTrigger): FlagTrigger.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FlagTrigger, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FlagTrigger;
  static deserializeBinaryFromReader(message: FlagTrigger, reader: jspb.BinaryReader): FlagTrigger;
}

export namespace FlagTrigger {
  export type AsObject = {
    id: string,
    featureId: string,
    environmentNamespace: string,
    type: FlagTrigger.TypeMap[keyof FlagTrigger.TypeMap],
    action: FlagTrigger.ActionMap[keyof FlagTrigger.ActionMap],
    description: string,
    triggerCount: number,
    lastTriggeredAt: number,
    uuid: string,
    disabled: boolean,
    createdAt: number,
    updatedAt: number,
    token: string,
  }

  export interface TypeMap {
    TYPE_UNKNOWN: 0;
    TYPE_WEBHOOK: 1;
  }

  export const Type: TypeMap;

  export interface ActionMap {
    ACTION_UNKNOWN: 0;
    ACTION_ON: 1;
    ACTION_OFF: 2;
  }

  export const Action: ActionMap;
}

