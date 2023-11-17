// package: bucketeer.autoops
// file: proto/autoops/flag_trigger.proto

import * as jspb from "google-protobuf";

export class FlagTrigger extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getType(): FlagTrigger.TiggerTypeMap[keyof FlagTrigger.TiggerTypeMap];
  setType(value: FlagTrigger.TiggerTypeMap[keyof FlagTrigger.TiggerTypeMap]): void;

  getAction(): TriggerActionMap[keyof TriggerActionMap];
  setAction(value: TriggerActionMap[keyof TriggerActionMap]): void;

  getDescription(): string;
  setDescription(value: string): void;

  getTriggerTimes(): number;
  setTriggerTimes(value: number): void;

  getLastTriggeredAt(): number;
  setLastTriggeredAt(value: number): void;

  getUuid(): string;
  setUuid(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

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
    type: FlagTrigger.TiggerTypeMap[keyof FlagTrigger.TiggerTypeMap],
    action: TriggerActionMap[keyof TriggerActionMap],
    description: string,
    triggerTimes: number,
    lastTriggeredAt: number,
    uuid: string,
    disabled: boolean,
    deleted: boolean,
    createdAt: number,
    updatedAt: number,
  }

  export interface TiggerTypeMap {
    UNKNOWN: 0;
    WEBHOOK: 1;
  }

  export const TiggerType: TiggerTypeMap;
}

export interface TriggerActionMap {
  UNKNOWN: 0;
  ON: 1;
  OFF: 2;
}

export const TriggerAction: TriggerActionMap;

