// package: bucketeer.autoops
// file: proto/autoops/command.proto

import * as jspb from 'google-protobuf';
import * as proto_autoops_auto_ops_rule_pb from '../../proto/autoops/auto_ops_rule_pb';
import * as proto_autoops_clause_pb from '../../proto/autoops/clause_pb';
import * as proto_autoops_progressive_rollout_pb from '../../proto/autoops/progressive_rollout_pb';

export class CreateAutoOpsRuleCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getOpsType(): proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  setOpsType(
    value: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap]
  ): void;

  clearOpsEventRateClausesList(): void;
  getOpsEventRateClausesList(): Array<proto_autoops_clause_pb.OpsEventRateClause>;
  setOpsEventRateClausesList(
    value: Array<proto_autoops_clause_pb.OpsEventRateClause>
  ): void;
  addOpsEventRateClauses(
    value?: proto_autoops_clause_pb.OpsEventRateClause,
    index?: number
  ): proto_autoops_clause_pb.OpsEventRateClause;

  clearDatetimeClausesList(): void;
  getDatetimeClausesList(): Array<proto_autoops_clause_pb.DatetimeClause>;
  setDatetimeClausesList(
    value: Array<proto_autoops_clause_pb.DatetimeClause>
  ): void;
  addDatetimeClauses(
    value?: proto_autoops_clause_pb.DatetimeClause,
    index?: number
  ): proto_autoops_clause_pb.DatetimeClause;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAutoOpsRuleCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAutoOpsRuleCommand
  ): CreateAutoOpsRuleCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAutoOpsRuleCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAutoOpsRuleCommand;
  static deserializeBinaryFromReader(
    message: CreateAutoOpsRuleCommand,
    reader: jspb.BinaryReader
  ): CreateAutoOpsRuleCommand;
}

export namespace CreateAutoOpsRuleCommand {
  export type AsObject = {
    featureId: string;
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
    opsEventRateClausesList: Array<proto_autoops_clause_pb.OpsEventRateClause.AsObject>;
    datetimeClausesList: Array<proto_autoops_clause_pb.DatetimeClause.AsObject>;
  };
}

export class ChangeAutoOpsRuleOpsTypeCommand extends jspb.Message {
  getOpsType(): proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  setOpsType(
    value: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAutoOpsRuleOpsTypeCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAutoOpsRuleOpsTypeCommand
  ): ChangeAutoOpsRuleOpsTypeCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAutoOpsRuleOpsTypeCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAutoOpsRuleOpsTypeCommand;
  static deserializeBinaryFromReader(
    message: ChangeAutoOpsRuleOpsTypeCommand,
    reader: jspb.BinaryReader
  ): ChangeAutoOpsRuleOpsTypeCommand;
}

export namespace ChangeAutoOpsRuleOpsTypeCommand {
  export type AsObject = {
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  };
}

export class DeleteAutoOpsRuleCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAutoOpsRuleCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteAutoOpsRuleCommand
  ): DeleteAutoOpsRuleCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteAutoOpsRuleCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAutoOpsRuleCommand;
  static deserializeBinaryFromReader(
    message: DeleteAutoOpsRuleCommand,
    reader: jspb.BinaryReader
  ): DeleteAutoOpsRuleCommand;
}

export namespace DeleteAutoOpsRuleCommand {
  export type AsObject = {};
}

export class StopAutoOpsRuleCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopAutoOpsRuleCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StopAutoOpsRuleCommand
  ): StopAutoOpsRuleCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StopAutoOpsRuleCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StopAutoOpsRuleCommand;
  static deserializeBinaryFromReader(
    message: StopAutoOpsRuleCommand,
    reader: jspb.BinaryReader
  ): StopAutoOpsRuleCommand;
}

export namespace StopAutoOpsRuleCommand {
  export type AsObject = {};
}

export class ChangeAutoOpsRuleTriggeredAtCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ChangeAutoOpsRuleTriggeredAtCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAutoOpsRuleTriggeredAtCommand
  ): ChangeAutoOpsRuleTriggeredAtCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAutoOpsRuleTriggeredAtCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ChangeAutoOpsRuleTriggeredAtCommand;
  static deserializeBinaryFromReader(
    message: ChangeAutoOpsRuleTriggeredAtCommand,
    reader: jspb.BinaryReader
  ): ChangeAutoOpsRuleTriggeredAtCommand;
}

export namespace ChangeAutoOpsRuleTriggeredAtCommand {
  export type AsObject = {};
}

export class ChangeAutoOpsStatusCommand extends jspb.Message {
  getStatus(): proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
  setStatus(
    value: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAutoOpsStatusCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeAutoOpsStatusCommand
  ): ChangeAutoOpsStatusCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeAutoOpsStatusCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAutoOpsStatusCommand;
  static deserializeBinaryFromReader(
    message: ChangeAutoOpsStatusCommand,
    reader: jspb.BinaryReader
  ): ChangeAutoOpsStatusCommand;
}

export namespace ChangeAutoOpsStatusCommand {
  export type AsObject = {
    status: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
  };
}

export class ExecuteAutoOpsRuleCommand extends jspb.Message {
  getClauseId(): string;
  setClauseId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteAutoOpsRuleCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExecuteAutoOpsRuleCommand
  ): ExecuteAutoOpsRuleCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExecuteAutoOpsRuleCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteAutoOpsRuleCommand;
  static deserializeBinaryFromReader(
    message: ExecuteAutoOpsRuleCommand,
    reader: jspb.BinaryReader
  ): ExecuteAutoOpsRuleCommand;
}

export namespace ExecuteAutoOpsRuleCommand {
  export type AsObject = {
    clauseId: string;
  };
}

export class AddOpsEventRateClauseCommand extends jspb.Message {
  hasOpsEventRateClause(): boolean;
  clearOpsEventRateClause(): void;
  getOpsEventRateClause():
    | proto_autoops_clause_pb.OpsEventRateClause
    | undefined;
  setOpsEventRateClause(
    value?: proto_autoops_clause_pb.OpsEventRateClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOpsEventRateClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AddOpsEventRateClauseCommand
  ): AddOpsEventRateClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AddOpsEventRateClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AddOpsEventRateClauseCommand;
  static deserializeBinaryFromReader(
    message: AddOpsEventRateClauseCommand,
    reader: jspb.BinaryReader
  ): AddOpsEventRateClauseCommand;
}

export namespace AddOpsEventRateClauseCommand {
  export type AsObject = {
    opsEventRateClause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject;
  };
}

export class ChangeOpsEventRateClauseCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasOpsEventRateClause(): boolean;
  clearOpsEventRateClause(): void;
  getOpsEventRateClause():
    | proto_autoops_clause_pb.OpsEventRateClause
    | undefined;
  setOpsEventRateClause(
    value?: proto_autoops_clause_pb.OpsEventRateClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeOpsEventRateClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeOpsEventRateClauseCommand
  ): ChangeOpsEventRateClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeOpsEventRateClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeOpsEventRateClauseCommand;
  static deserializeBinaryFromReader(
    message: ChangeOpsEventRateClauseCommand,
    reader: jspb.BinaryReader
  ): ChangeOpsEventRateClauseCommand;
}

export namespace ChangeOpsEventRateClauseCommand {
  export type AsObject = {
    id: string;
    opsEventRateClause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject;
  };
}

export class DeleteClauseCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteClauseCommand
  ): DeleteClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteClauseCommand;
  static deserializeBinaryFromReader(
    message: DeleteClauseCommand,
    reader: jspb.BinaryReader
  ): DeleteClauseCommand;
}

export namespace DeleteClauseCommand {
  export type AsObject = {
    id: string;
  };
}

export class AddDatetimeClauseCommand extends jspb.Message {
  hasDatetimeClause(): boolean;
  clearDatetimeClause(): void;
  getDatetimeClause(): proto_autoops_clause_pb.DatetimeClause | undefined;
  setDatetimeClause(value?: proto_autoops_clause_pb.DatetimeClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDatetimeClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AddDatetimeClauseCommand
  ): AddDatetimeClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AddDatetimeClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AddDatetimeClauseCommand;
  static deserializeBinaryFromReader(
    message: AddDatetimeClauseCommand,
    reader: jspb.BinaryReader
  ): AddDatetimeClauseCommand;
}

export namespace AddDatetimeClauseCommand {
  export type AsObject = {
    datetimeClause?: proto_autoops_clause_pb.DatetimeClause.AsObject;
  };
}

export class ChangeDatetimeClauseCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasDatetimeClause(): boolean;
  clearDatetimeClause(): void;
  getDatetimeClause(): proto_autoops_clause_pb.DatetimeClause | undefined;
  setDatetimeClause(value?: proto_autoops_clause_pb.DatetimeClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeDatetimeClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeDatetimeClauseCommand
  ): ChangeDatetimeClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeDatetimeClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDatetimeClauseCommand;
  static deserializeBinaryFromReader(
    message: ChangeDatetimeClauseCommand,
    reader: jspb.BinaryReader
  ): ChangeDatetimeClauseCommand;
}

export namespace ChangeDatetimeClauseCommand {
  export type AsObject = {
    id: string;
    datetimeClause?: proto_autoops_clause_pb.DatetimeClause.AsObject;
  };
}

export class CreateProgressiveRolloutCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  hasProgressiveRolloutManualScheduleClause(): boolean;
  clearProgressiveRolloutManualScheduleClause(): void;
  getProgressiveRolloutManualScheduleClause():
    | proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause
    | undefined;
  setProgressiveRolloutManualScheduleClause(
    value?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause
  ): void;

  hasProgressiveRolloutTemplateScheduleClause(): boolean;
  clearProgressiveRolloutTemplateScheduleClause(): void;
  getProgressiveRolloutTemplateScheduleClause():
    | proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause
    | undefined;
  setProgressiveRolloutTemplateScheduleClause(
    value?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProgressiveRolloutCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateProgressiveRolloutCommand
  ): CreateProgressiveRolloutCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateProgressiveRolloutCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateProgressiveRolloutCommand;
  static deserializeBinaryFromReader(
    message: CreateProgressiveRolloutCommand,
    reader: jspb.BinaryReader
  ): CreateProgressiveRolloutCommand;
}

export namespace CreateProgressiveRolloutCommand {
  export type AsObject = {
    featureId: string;
    progressiveRolloutManualScheduleClause?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause.AsObject;
    progressiveRolloutTemplateScheduleClause?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause.AsObject;
  };
}

export class StopProgressiveRolloutCommand extends jspb.Message {
  getStoppedBy(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap];
  setStoppedBy(
    value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopProgressiveRolloutCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StopProgressiveRolloutCommand
  ): StopProgressiveRolloutCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StopProgressiveRolloutCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StopProgressiveRolloutCommand;
  static deserializeBinaryFromReader(
    message: StopProgressiveRolloutCommand,
    reader: jspb.BinaryReader
  ): StopProgressiveRolloutCommand;
}

export namespace StopProgressiveRolloutCommand {
  export type AsObject = {
    stoppedBy: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap];
  };
}

export class DeleteProgressiveRolloutCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProgressiveRolloutCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteProgressiveRolloutCommand
  ): DeleteProgressiveRolloutCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteProgressiveRolloutCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProgressiveRolloutCommand;
  static deserializeBinaryFromReader(
    message: DeleteProgressiveRolloutCommand,
    reader: jspb.BinaryReader
  ): DeleteProgressiveRolloutCommand;
}

export namespace DeleteProgressiveRolloutCommand {
  export type AsObject = {};
}

export class AddProgressiveRolloutManualScheduleClauseCommand extends jspb.Message {
  hasClause(): boolean;
  clearClause(): void;
  getClause():
    | proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause
    | undefined;
  setClause(
    value?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AddProgressiveRolloutManualScheduleClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AddProgressiveRolloutManualScheduleClauseCommand
  ): AddProgressiveRolloutManualScheduleClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AddProgressiveRolloutManualScheduleClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AddProgressiveRolloutManualScheduleClauseCommand;
  static deserializeBinaryFromReader(
    message: AddProgressiveRolloutManualScheduleClauseCommand,
    reader: jspb.BinaryReader
  ): AddProgressiveRolloutManualScheduleClauseCommand;
}

export namespace AddProgressiveRolloutManualScheduleClauseCommand {
  export type AsObject = {
    clause?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause.AsObject;
  };
}

export class AddProgressiveRolloutTemplateScheduleClauseCommand extends jspb.Message {
  hasClause(): boolean;
  clearClause(): void;
  getClause():
    | proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause
    | undefined;
  setClause(
    value?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AddProgressiveRolloutTemplateScheduleClauseCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AddProgressiveRolloutTemplateScheduleClauseCommand
  ): AddProgressiveRolloutTemplateScheduleClauseCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AddProgressiveRolloutTemplateScheduleClauseCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AddProgressiveRolloutTemplateScheduleClauseCommand;
  static deserializeBinaryFromReader(
    message: AddProgressiveRolloutTemplateScheduleClauseCommand,
    reader: jspb.BinaryReader
  ): AddProgressiveRolloutTemplateScheduleClauseCommand;
}

export namespace AddProgressiveRolloutTemplateScheduleClauseCommand {
  export type AsObject = {
    clause?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause.AsObject;
  };
}

export class ChangeProgressiveRolloutScheduleTriggeredAtCommand extends jspb.Message {
  getScheduleId(): string;
  setScheduleId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ChangeProgressiveRolloutScheduleTriggeredAtCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ChangeProgressiveRolloutScheduleTriggeredAtCommand
  ): ChangeProgressiveRolloutScheduleTriggeredAtCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ChangeProgressiveRolloutScheduleTriggeredAtCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ChangeProgressiveRolloutScheduleTriggeredAtCommand;
  static deserializeBinaryFromReader(
    message: ChangeProgressiveRolloutScheduleTriggeredAtCommand,
    reader: jspb.BinaryReader
  ): ChangeProgressiveRolloutScheduleTriggeredAtCommand;
}

export namespace ChangeProgressiveRolloutScheduleTriggeredAtCommand {
  export type AsObject = {
    scheduleId: string;
  };
}
