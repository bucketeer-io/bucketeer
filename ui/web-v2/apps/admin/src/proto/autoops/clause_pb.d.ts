// package: bucketeer.autoops
// file: proto/autoops/clause.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";

export class Clause extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasClause(): boolean;
  clearClause(): void;
  getClause(): google_protobuf_any_pb.Any | undefined;
  setClause(value?: google_protobuf_any_pb.Any): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Clause.AsObject;
  static toObject(includeInstance: boolean, msg: Clause): Clause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Clause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Clause;
  static deserializeBinaryFromReader(message: Clause, reader: jspb.BinaryReader): Clause;
}

export namespace Clause {
  export type AsObject = {
    id: string,
    clause?: google_protobuf_any_pb.Any.AsObject,
  }
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
  setOperator(value: OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpsEventRateClause.AsObject;
  static toObject(includeInstance: boolean, msg: OpsEventRateClause): OpsEventRateClause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpsEventRateClause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpsEventRateClause;
  static deserializeBinaryFromReader(message: OpsEventRateClause, reader: jspb.BinaryReader): OpsEventRateClause;
}

export namespace OpsEventRateClause {
  export type AsObject = {
    variationId: string,
    goalId: string,
    minCount: number,
    threadsholdRate: number,
    operator: OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap],
  }

  export interface OperatorMap {
    GREATER_OR_EQUAL: 0;
    LESS_OR_EQUAL: 1;
  }

  export const Operator: OperatorMap;
}

export class DatetimeClause extends jspb.Message {
  getTime(): number;
  setTime(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatetimeClause.AsObject;
  static toObject(includeInstance: boolean, msg: DatetimeClause): DatetimeClause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DatetimeClause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DatetimeClause;
  static deserializeBinaryFromReader(message: DatetimeClause, reader: jspb.BinaryReader): DatetimeClause;
}

export namespace DatetimeClause {
  export type AsObject = {
    time: number,
  }
}

export class WebhookClause extends jspb.Message {
  getWebhookId(): string;
  setWebhookId(value: string): void;

  clearConditionsList(): void;
  getConditionsList(): Array<WebhookClause.Condition>;
  setConditionsList(value: Array<WebhookClause.Condition>): void;
  addConditions(value?: WebhookClause.Condition, index?: number): WebhookClause.Condition;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebhookClause.AsObject;
  static toObject(includeInstance: boolean, msg: WebhookClause): WebhookClause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WebhookClause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebhookClause;
  static deserializeBinaryFromReader(message: WebhookClause, reader: jspb.BinaryReader): WebhookClause;
}

export namespace WebhookClause {
  export type AsObject = {
    webhookId: string,
    conditionsList: Array<WebhookClause.Condition.AsObject>,
  }

  export class Condition extends jspb.Message {
    getFilter(): string;
    setFilter(value: string): void;

    getValue(): string;
    setValue(value: string): void;

    getOperator(): WebhookClause.Condition.OperatorMap[keyof WebhookClause.Condition.OperatorMap];
    setOperator(value: WebhookClause.Condition.OperatorMap[keyof WebhookClause.Condition.OperatorMap]): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Condition.AsObject;
    static toObject(includeInstance: boolean, msg: Condition): Condition.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Condition, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Condition;
    static deserializeBinaryFromReader(message: Condition, reader: jspb.BinaryReader): Condition;
  }

  export namespace Condition {
    export type AsObject = {
      filter: string,
      value: string,
      operator: WebhookClause.Condition.OperatorMap[keyof WebhookClause.Condition.OperatorMap],
    }

    export interface OperatorMap {
      EQUAL: 0;
      NOT_EQUAL: 1;
      MORE_THAN: 2;
      MORE_THAN_OR_EQUAL: 3;
      LESS_THAN: 4;
      LESS_THAN_OR_EQUAL: 5;
    }

    export const Operator: OperatorMap;
  }
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
  static toObject(includeInstance: boolean, msg: ProgressiveRolloutSchedule): ProgressiveRolloutSchedule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProgressiveRolloutSchedule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutSchedule;
  static deserializeBinaryFromReader(message: ProgressiveRolloutSchedule, reader: jspb.BinaryReader): ProgressiveRolloutSchedule;
}

export namespace ProgressiveRolloutSchedule {
  export type AsObject = {
    scheduleId: string,
    executeAt: number,
    weight: number,
    triggeredAt: number,
  }
}

export class ProgressiveRolloutManualScheduleClause extends jspb.Message {
  clearSchedulesList(): void;
  getSchedulesList(): Array<ProgressiveRolloutSchedule>;
  setSchedulesList(value: Array<ProgressiveRolloutSchedule>): void;
  addSchedules(value?: ProgressiveRolloutSchedule, index?: number): ProgressiveRolloutSchedule;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRolloutManualScheduleClause.AsObject;
  static toObject(includeInstance: boolean, msg: ProgressiveRolloutManualScheduleClause): ProgressiveRolloutManualScheduleClause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProgressiveRolloutManualScheduleClause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutManualScheduleClause;
  static deserializeBinaryFromReader(message: ProgressiveRolloutManualScheduleClause, reader: jspb.BinaryReader): ProgressiveRolloutManualScheduleClause;
}

export namespace ProgressiveRolloutManualScheduleClause {
  export type AsObject = {
    schedulesList: Array<ProgressiveRolloutSchedule.AsObject>,
    variationId: string,
  }
}

export class ProgressiveRolloutTemplateScheduleClause extends jspb.Message {
  clearSchedulesList(): void;
  getSchedulesList(): Array<ProgressiveRolloutSchedule>;
  setSchedulesList(value: Array<ProgressiveRolloutSchedule>): void;
  addSchedules(value?: ProgressiveRolloutSchedule, index?: number): ProgressiveRolloutSchedule;

  getInterval(): ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap];
  setInterval(value: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap]): void;

  getIncrements(): number;
  setIncrements(value: number): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRolloutTemplateScheduleClause.AsObject;
  static toObject(includeInstance: boolean, msg: ProgressiveRolloutTemplateScheduleClause): ProgressiveRolloutTemplateScheduleClause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProgressiveRolloutTemplateScheduleClause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutTemplateScheduleClause;
  static deserializeBinaryFromReader(message: ProgressiveRolloutTemplateScheduleClause, reader: jspb.BinaryReader): ProgressiveRolloutTemplateScheduleClause;
}

export namespace ProgressiveRolloutTemplateScheduleClause {
  export type AsObject = {
    schedulesList: Array<ProgressiveRolloutSchedule.AsObject>,
    interval: ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap],
    increments: number,
    variationId: string,
  }

  export interface IntervalMap {
    UNKNOWN: 0;
    HOURLY: 1;
    DAILY: 2;
    WEEKLY: 3;
  }

  export const Interval: IntervalMap;
}

