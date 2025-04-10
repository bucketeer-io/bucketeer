// package: bucketeer.autoops
// file: proto/autoops/clause.proto

import * as jspb from 'google-protobuf';
import * as google_protobuf_any_pb from 'google-protobuf/google/protobuf/any_pb';

export class Clause extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasClause(): boolean;
  clearClause(): void;
  getClause(): google_protobuf_any_pb.Any | undefined;
  setClause(value?: google_protobuf_any_pb.Any): void;

  getActionType(): ActionTypeMap[keyof ActionTypeMap];
  setActionType(value: ActionTypeMap[keyof ActionTypeMap]): void;

  getExecutedAt(): number;
  setExecutedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Clause.AsObject;
  static toObject(includeInstance: boolean, msg: Clause): Clause.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Clause,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Clause;
  static deserializeBinaryFromReader(
    message: Clause,
    reader: jspb.BinaryReader
  ): Clause;
}

export namespace Clause {
  export type AsObject = {
    id: string;
    clause?: google_protobuf_any_pb.Any.AsObject;
    actionType: ActionTypeMap[keyof ActionTypeMap];
    executedAt: number;
  };
}

export class OpsEventRateClause extends jspb.Message {
  getVariationId(): string;
  setVariationId(value: string): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getMinCount(): number;
  setMinCount(value: number): void;

  getThreadsholdRate(): number;
  setThreadsholdRate(value: number): void;

  getOperator(): OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap];
  setOperator(
    value: OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap]
  ): void;

  getActionType(): ActionTypeMap[keyof ActionTypeMap];
  setActionType(value: ActionTypeMap[keyof ActionTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpsEventRateClause.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OpsEventRateClause
  ): OpsEventRateClause.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OpsEventRateClause,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OpsEventRateClause;
  static deserializeBinaryFromReader(
    message: OpsEventRateClause,
    reader: jspb.BinaryReader
  ): OpsEventRateClause;
}

export namespace OpsEventRateClause {
  export type AsObject = {
    variationId: string;
    goalId: string;
    minCount: number;
    threadsholdRate: number;
    operator: OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap];
    actionType: ActionTypeMap[keyof ActionTypeMap];
  };

  export interface OperatorMap {
    GREATER_OR_EQUAL: 0;
    LESS_OR_EQUAL: 1;
  }

  export const Operator: OperatorMap;
}

export class DatetimeClause extends jspb.Message {
  getTime(): number;
  setTime(value: number): void;

  getActionType(): ActionTypeMap[keyof ActionTypeMap];
  setActionType(value: ActionTypeMap[keyof ActionTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatetimeClause.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DatetimeClause
  ): DatetimeClause.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DatetimeClause,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DatetimeClause;
  static deserializeBinaryFromReader(
    message: DatetimeClause,
    reader: jspb.BinaryReader
  ): DatetimeClause;
}

export namespace DatetimeClause {
  export type AsObject = {
    time: number;
    actionType: ActionTypeMap[keyof ActionTypeMap];
  };
}

export class ProgressiveRolloutSchedule extends jspb.Message {
  getScheduleId(): string;
  setScheduleId(value: string): void;

  getExecuteAt(): number;
  setExecuteAt(value: number): void;

  getWeight(): number;
  setWeight(value: number): void;

  getTriggeredAt(): number;
  setTriggeredAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRolloutSchedule.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutSchedule
  ): ProgressiveRolloutSchedule.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutSchedule,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutSchedule;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutSchedule,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutSchedule;
}

export namespace ProgressiveRolloutSchedule {
  export type AsObject = {
    scheduleId: string;
    executeAt: number;
    weight: number;
    triggeredAt: number;
  };
}

export class ProgressiveRolloutManualScheduleClause extends jspb.Message {
  clearSchedulesList(): void;
  getSchedulesList(): Array<ProgressiveRolloutSchedule>;
  setSchedulesList(value: Array<ProgressiveRolloutSchedule>): void;
  addSchedules(
    value?: ProgressiveRolloutSchedule,
    index?: number
  ): ProgressiveRolloutSchedule;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ProgressiveRolloutManualScheduleClause.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutManualScheduleClause
  ): ProgressiveRolloutManualScheduleClause.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutManualScheduleClause,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ProgressiveRolloutManualScheduleClause;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutManualScheduleClause,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutManualScheduleClause;
}

export namespace ProgressiveRolloutManualScheduleClause {
  export type AsObject = {
    schedulesList: Array<ProgressiveRolloutSchedule.AsObject>;
    variationId: string;
  };
}

export class ProgressiveRolloutTemplateScheduleClause extends jspb.Message {
  clearSchedulesList(): void;
  getSchedulesList(): Array<ProgressiveRolloutSchedule>;
  setSchedulesList(value: Array<ProgressiveRolloutSchedule>): void;
  addSchedules(
    value?: ProgressiveRolloutSchedule,
    index?: number
  ): ProgressiveRolloutSchedule;

  getInterval(): ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap];
  setInterval(
    value: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap]
  ): void;

  getIncrements(): number;
  setIncrements(value: number): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ProgressiveRolloutTemplateScheduleClause.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutTemplateScheduleClause
  ): ProgressiveRolloutTemplateScheduleClause.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutTemplateScheduleClause,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ProgressiveRolloutTemplateScheduleClause;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutTemplateScheduleClause,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutTemplateScheduleClause;
}

export namespace ProgressiveRolloutTemplateScheduleClause {
  export type AsObject = {
    schedulesList: Array<ProgressiveRolloutSchedule.AsObject>;
    interval: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap];
    increments: number;
    variationId: string;
  };

  export interface IntervalMap {
    UNKNOWN: 0;
    HOURLY: 1;
    DAILY: 2;
    WEEKLY: 3;
  }

  export const Interval: IntervalMap;
}

export interface ActionTypeMap {
  UNKNOWN: 0;
  ENABLE: 1;
  DISABLE: 2;
}

export const ActionType: ActionTypeMap;
