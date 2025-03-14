// package: bucketeer.feature
// file: proto/feature/scheduled_update.proto

import * as jspb from 'google-protobuf';

export class ScheduledChange extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getChangeType(): ScheduledChange.ChangeTypeMap[keyof ScheduledChange.ChangeTypeMap];
  setChangeType(
    value: ScheduledChange.ChangeTypeMap[keyof ScheduledChange.ChangeTypeMap]
  ): void;

  getFieldType(): ScheduledChange.FieldTypeMap[keyof ScheduledChange.FieldTypeMap];
  setFieldType(
    value: ScheduledChange.FieldTypeMap[keyof ScheduledChange.FieldTypeMap]
  ): void;

  getFieldValue(): string;
  setFieldValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ScheduledChange.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ScheduledChange
  ): ScheduledChange.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ScheduledChange,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ScheduledChange;
  static deserializeBinaryFromReader(
    message: ScheduledChange,
    reader: jspb.BinaryReader
  ): ScheduledChange;
}

export namespace ScheduledChange {
  export type AsObject = {
    id: string;
    changeType: ScheduledChange.ChangeTypeMap[keyof ScheduledChange.ChangeTypeMap];
    fieldType: ScheduledChange.FieldTypeMap[keyof ScheduledChange.FieldTypeMap];
    fieldValue: string;
  };

  export interface FieldTypeMap {
    UNSPECIFIED: 0;
    PREREQUISITES: 1;
    TARGETS: 2;
    RULES: 3;
    DEFAULT_STRATEGY: 4;
    OFF_VARIATION: 5;
    VARIATIONS: 6;
  }

  export const FieldType: FieldTypeMap;

  export interface ChangeTypeMap {
    CHANGE_UNSPECIFIED: 0;
    CHANGE_CREATE: 1;
    CHANGE_UPDATE: 2;
    CHANGE_DELETE: 3;
  }

  export const ChangeType: ChangeTypeMap;
}

export class ScheduledFlagUpdate extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getScheduledAt(): number;
  setScheduledAt(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  clearChangesList(): void;
  getChangesList(): Array<ScheduledChange>;
  setChangesList(value: Array<ScheduledChange>): void;
  addChanges(value?: ScheduledChange, index?: number): ScheduledChange;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ScheduledFlagUpdate.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ScheduledFlagUpdate
  ): ScheduledFlagUpdate.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ScheduledFlagUpdate,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ScheduledFlagUpdate;
  static deserializeBinaryFromReader(
    message: ScheduledFlagUpdate,
    reader: jspb.BinaryReader
  ): ScheduledFlagUpdate;
}

export namespace ScheduledFlagUpdate {
  export type AsObject = {
    id: string;
    featureId: string;
    environmentId: string;
    scheduledAt: number;
    createdAt: number;
    updatedAt: number;
    changesList: Array<ScheduledChange.AsObject>;
  };
}
