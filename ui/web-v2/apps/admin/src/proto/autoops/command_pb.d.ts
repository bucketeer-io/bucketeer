// package: bucketeer.autoops
// file: proto/autoops/command.proto

import * as jspb from "google-protobuf";
import * as proto_autoops_auto_ops_rule_pb from "../../proto/autoops/auto_ops_rule_pb";
import * as proto_autoops_clause_pb from "../../proto/autoops/clause_pb";
import * as proto_autoops_flag_trigger_pb from "../../proto/autoops/flag_trigger_pb";

export class CreateAutoOpsRuleCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getOpsType(): proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  setOpsType(value: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap]): void;

  clearOpsEventRateClausesList(): void;
  getOpsEventRateClausesList(): Array<proto_autoops_clause_pb.OpsEventRateClause>;
  setOpsEventRateClausesList(value: Array<proto_autoops_clause_pb.OpsEventRateClause>): void;
  addOpsEventRateClauses(value?: proto_autoops_clause_pb.OpsEventRateClause, index?: number): proto_autoops_clause_pb.OpsEventRateClause;

  clearDatetimeClausesList(): void;
  getDatetimeClausesList(): Array<proto_autoops_clause_pb.DatetimeClause>;
  setDatetimeClausesList(value: Array<proto_autoops_clause_pb.DatetimeClause>): void;
  addDatetimeClauses(value?: proto_autoops_clause_pb.DatetimeClause, index?: number): proto_autoops_clause_pb.DatetimeClause;

  clearWebhookClausesList(): void;
  getWebhookClausesList(): Array<proto_autoops_clause_pb.WebhookClause>;
  setWebhookClausesList(value: Array<proto_autoops_clause_pb.WebhookClause>): void;
  addWebhookClauses(value?: proto_autoops_clause_pb.WebhookClause, index?: number): proto_autoops_clause_pb.WebhookClause;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAutoOpsRuleCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAutoOpsRuleCommand): CreateAutoOpsRuleCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAutoOpsRuleCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAutoOpsRuleCommand;
  static deserializeBinaryFromReader(message: CreateAutoOpsRuleCommand, reader: jspb.BinaryReader): CreateAutoOpsRuleCommand;
}

export namespace CreateAutoOpsRuleCommand {
  export type AsObject = {
    featureId: string,
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap],
    opsEventRateClausesList: Array<proto_autoops_clause_pb.OpsEventRateClause.AsObject>,
    datetimeClausesList: Array<proto_autoops_clause_pb.DatetimeClause.AsObject>,
    webhookClausesList: Array<proto_autoops_clause_pb.WebhookClause.AsObject>,
  }
}

export class ChangeAutoOpsRuleOpsTypeCommand extends jspb.Message {
  getOpsType(): proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  setOpsType(value: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAutoOpsRuleOpsTypeCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAutoOpsRuleOpsTypeCommand): ChangeAutoOpsRuleOpsTypeCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAutoOpsRuleOpsTypeCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAutoOpsRuleOpsTypeCommand;
  static deserializeBinaryFromReader(message: ChangeAutoOpsRuleOpsTypeCommand, reader: jspb.BinaryReader): ChangeAutoOpsRuleOpsTypeCommand;
}

export namespace ChangeAutoOpsRuleOpsTypeCommand {
  export type AsObject = {
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap],
  }
}

export class DeleteAutoOpsRuleCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAutoOpsRuleCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAutoOpsRuleCommand): DeleteAutoOpsRuleCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAutoOpsRuleCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAutoOpsRuleCommand;
  static deserializeBinaryFromReader(message: DeleteAutoOpsRuleCommand, reader: jspb.BinaryReader): DeleteAutoOpsRuleCommand;
}

export namespace DeleteAutoOpsRuleCommand {
  export type AsObject = {
  }
}

export class ChangeAutoOpsRuleTriggeredAtCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeAutoOpsRuleTriggeredAtCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeAutoOpsRuleTriggeredAtCommand): ChangeAutoOpsRuleTriggeredAtCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeAutoOpsRuleTriggeredAtCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeAutoOpsRuleTriggeredAtCommand;
  static deserializeBinaryFromReader(message: ChangeAutoOpsRuleTriggeredAtCommand, reader: jspb.BinaryReader): ChangeAutoOpsRuleTriggeredAtCommand;
}

export namespace ChangeAutoOpsRuleTriggeredAtCommand {
  export type AsObject = {
  }
}

export class AddOpsEventRateClauseCommand extends jspb.Message {
  hasOpsEventRateClause(): boolean;
  clearOpsEventRateClause(): void;
  getOpsEventRateClause(): proto_autoops_clause_pb.OpsEventRateClause | undefined;
  setOpsEventRateClause(value?: proto_autoops_clause_pb.OpsEventRateClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOpsEventRateClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddOpsEventRateClauseCommand): AddOpsEventRateClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddOpsEventRateClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOpsEventRateClauseCommand;
  static deserializeBinaryFromReader(message: AddOpsEventRateClauseCommand, reader: jspb.BinaryReader): AddOpsEventRateClauseCommand;
}

export namespace AddOpsEventRateClauseCommand {
  export type AsObject = {
    opsEventRateClause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject,
  }
}

export class ChangeOpsEventRateClauseCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasOpsEventRateClause(): boolean;
  clearOpsEventRateClause(): void;
  getOpsEventRateClause(): proto_autoops_clause_pb.OpsEventRateClause | undefined;
  setOpsEventRateClause(value?: proto_autoops_clause_pb.OpsEventRateClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeOpsEventRateClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeOpsEventRateClauseCommand): ChangeOpsEventRateClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeOpsEventRateClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeOpsEventRateClauseCommand;
  static deserializeBinaryFromReader(message: ChangeOpsEventRateClauseCommand, reader: jspb.BinaryReader): ChangeOpsEventRateClauseCommand;
}

export namespace ChangeOpsEventRateClauseCommand {
  export type AsObject = {
    id: string,
    opsEventRateClause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject,
  }
}

export class DeleteClauseCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteClauseCommand): DeleteClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteClauseCommand;
  static deserializeBinaryFromReader(message: DeleteClauseCommand, reader: jspb.BinaryReader): DeleteClauseCommand;
}

export namespace DeleteClauseCommand {
  export type AsObject = {
    id: string,
  }
}

export class AddDatetimeClauseCommand extends jspb.Message {
  hasDatetimeClause(): boolean;
  clearDatetimeClause(): void;
  getDatetimeClause(): proto_autoops_clause_pb.DatetimeClause | undefined;
  setDatetimeClause(value?: proto_autoops_clause_pb.DatetimeClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDatetimeClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddDatetimeClauseCommand): AddDatetimeClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddDatetimeClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDatetimeClauseCommand;
  static deserializeBinaryFromReader(message: AddDatetimeClauseCommand, reader: jspb.BinaryReader): AddDatetimeClauseCommand;
}

export namespace AddDatetimeClauseCommand {
  export type AsObject = {
    datetimeClause?: proto_autoops_clause_pb.DatetimeClause.AsObject,
  }
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
  static toObject(includeInstance: boolean, msg: ChangeDatetimeClauseCommand): ChangeDatetimeClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeDatetimeClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDatetimeClauseCommand;
  static deserializeBinaryFromReader(message: ChangeDatetimeClauseCommand, reader: jspb.BinaryReader): ChangeDatetimeClauseCommand;
}

export namespace ChangeDatetimeClauseCommand {
  export type AsObject = {
    id: string,
    datetimeClause?: proto_autoops_clause_pb.DatetimeClause.AsObject,
  }
}

export class CreateWebhookCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWebhookCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWebhookCommand): CreateWebhookCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateWebhookCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWebhookCommand;
  static deserializeBinaryFromReader(message: CreateWebhookCommand, reader: jspb.BinaryReader): CreateWebhookCommand;
}

export namespace CreateWebhookCommand {
  export type AsObject = {
    name: string,
    description: string,
  }
}

export class ChangeWebhookNameCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeWebhookNameCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeWebhookNameCommand): ChangeWebhookNameCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeWebhookNameCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeWebhookNameCommand;
  static deserializeBinaryFromReader(message: ChangeWebhookNameCommand, reader: jspb.BinaryReader): ChangeWebhookNameCommand;
}

export namespace ChangeWebhookNameCommand {
  export type AsObject = {
    name: string,
  }
}

export class ChangeWebhookDescriptionCommand extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeWebhookDescriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeWebhookDescriptionCommand): ChangeWebhookDescriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeWebhookDescriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeWebhookDescriptionCommand;
  static deserializeBinaryFromReader(message: ChangeWebhookDescriptionCommand, reader: jspb.BinaryReader): ChangeWebhookDescriptionCommand;
}

export namespace ChangeWebhookDescriptionCommand {
  export type AsObject = {
    description: string,
  }
}

export class DeleteWebhookCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteWebhookCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteWebhookCommand): DeleteWebhookCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteWebhookCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteWebhookCommand;
  static deserializeBinaryFromReader(message: DeleteWebhookCommand, reader: jspb.BinaryReader): DeleteWebhookCommand;
}

export namespace DeleteWebhookCommand {
  export type AsObject = {
  }
}

export class AddWebhookClauseCommand extends jspb.Message {
  hasWebhookClause(): boolean;
  clearWebhookClause(): void;
  getWebhookClause(): proto_autoops_clause_pb.WebhookClause | undefined;
  setWebhookClause(value?: proto_autoops_clause_pb.WebhookClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddWebhookClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddWebhookClauseCommand): AddWebhookClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddWebhookClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddWebhookClauseCommand;
  static deserializeBinaryFromReader(message: AddWebhookClauseCommand, reader: jspb.BinaryReader): AddWebhookClauseCommand;
}

export namespace AddWebhookClauseCommand {
  export type AsObject = {
    webhookClause?: proto_autoops_clause_pb.WebhookClause.AsObject,
  }
}

export class ChangeWebhookClauseCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasWebhookClause(): boolean;
  clearWebhookClause(): void;
  getWebhookClause(): proto_autoops_clause_pb.WebhookClause | undefined;
  setWebhookClause(value?: proto_autoops_clause_pb.WebhookClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeWebhookClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeWebhookClauseCommand): ChangeWebhookClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeWebhookClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeWebhookClauseCommand;
  static deserializeBinaryFromReader(message: ChangeWebhookClauseCommand, reader: jspb.BinaryReader): ChangeWebhookClauseCommand;
}

export namespace ChangeWebhookClauseCommand {
  export type AsObject = {
    id: string,
    webhookClause?: proto_autoops_clause_pb.WebhookClause.AsObject,
  }
}

export class CreateProgressiveRolloutCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  hasProgressiveRolloutManualScheduleClause(): boolean;
  clearProgressiveRolloutManualScheduleClause(): void;
  getProgressiveRolloutManualScheduleClause(): proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause | undefined;
  setProgressiveRolloutManualScheduleClause(value?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause): void;

  hasProgressiveRolloutTemplateScheduleClause(): boolean;
  clearProgressiveRolloutTemplateScheduleClause(): void;
  getProgressiveRolloutTemplateScheduleClause(): proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause | undefined;
  setProgressiveRolloutTemplateScheduleClause(value?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProgressiveRolloutCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateProgressiveRolloutCommand): CreateProgressiveRolloutCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateProgressiveRolloutCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateProgressiveRolloutCommand;
  static deserializeBinaryFromReader(message: CreateProgressiveRolloutCommand, reader: jspb.BinaryReader): CreateProgressiveRolloutCommand;
}

export namespace CreateProgressiveRolloutCommand {
  export type AsObject = {
    featureId: string,
    progressiveRolloutManualScheduleClause?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause.AsObject,
    progressiveRolloutTemplateScheduleClause?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause.AsObject,
  }
}

export class DeleteProgressiveRolloutCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProgressiveRolloutCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProgressiveRolloutCommand): DeleteProgressiveRolloutCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteProgressiveRolloutCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProgressiveRolloutCommand;
  static deserializeBinaryFromReader(message: DeleteProgressiveRolloutCommand, reader: jspb.BinaryReader): DeleteProgressiveRolloutCommand;
}

export namespace DeleteProgressiveRolloutCommand {
  export type AsObject = {
  }
}

export class AddProgressiveRolloutManualScheduleClauseCommand extends jspb.Message {
  hasClause(): boolean;
  clearClause(): void;
  getClause(): proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause | undefined;
  setClause(value?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProgressiveRolloutManualScheduleClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddProgressiveRolloutManualScheduleClauseCommand): AddProgressiveRolloutManualScheduleClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddProgressiveRolloutManualScheduleClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProgressiveRolloutManualScheduleClauseCommand;
  static deserializeBinaryFromReader(message: AddProgressiveRolloutManualScheduleClauseCommand, reader: jspb.BinaryReader): AddProgressiveRolloutManualScheduleClauseCommand;
}

export namespace AddProgressiveRolloutManualScheduleClauseCommand {
  export type AsObject = {
    clause?: proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause.AsObject,
  }
}

export class AddProgressiveRolloutTemplateScheduleClauseCommand extends jspb.Message {
  hasClause(): boolean;
  clearClause(): void;
  getClause(): proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause | undefined;
  setClause(value?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProgressiveRolloutTemplateScheduleClauseCommand.AsObject;
  static toObject(includeInstance: boolean, msg: AddProgressiveRolloutTemplateScheduleClauseCommand): AddProgressiveRolloutTemplateScheduleClauseCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddProgressiveRolloutTemplateScheduleClauseCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProgressiveRolloutTemplateScheduleClauseCommand;
  static deserializeBinaryFromReader(message: AddProgressiveRolloutTemplateScheduleClauseCommand, reader: jspb.BinaryReader): AddProgressiveRolloutTemplateScheduleClauseCommand;
}

export namespace AddProgressiveRolloutTemplateScheduleClauseCommand {
  export type AsObject = {
    clause?: proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause.AsObject,
  }
}

export class ChangeProgressiveRolloutScheduleTriggeredAtCommand extends jspb.Message {
  getScheduleId(): string;
  setScheduleId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeProgressiveRolloutScheduleTriggeredAtCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeProgressiveRolloutScheduleTriggeredAtCommand): ChangeProgressiveRolloutScheduleTriggeredAtCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeProgressiveRolloutScheduleTriggeredAtCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeProgressiveRolloutScheduleTriggeredAtCommand;
  static deserializeBinaryFromReader(message: ChangeProgressiveRolloutScheduleTriggeredAtCommand, reader: jspb.BinaryReader): ChangeProgressiveRolloutScheduleTriggeredAtCommand;
}

export namespace ChangeProgressiveRolloutScheduleTriggeredAtCommand {
  export type AsObject = {
    scheduleId: string,
  }
}

export class CreateFlagTriggerCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getType(): proto_autoops_flag_trigger_pb.FlagTrigger.TypeMap[keyof proto_autoops_flag_trigger_pb.FlagTrigger.TypeMap];
  setType(value: proto_autoops_flag_trigger_pb.FlagTrigger.TypeMap[keyof proto_autoops_flag_trigger_pb.FlagTrigger.TypeMap]): void;

  getAction(): proto_autoops_flag_trigger_pb.FlagTrigger.ActionMap[keyof proto_autoops_flag_trigger_pb.FlagTrigger.ActionMap];
  setAction(value: proto_autoops_flag_trigger_pb.FlagTrigger.ActionMap[keyof proto_autoops_flag_trigger_pb.FlagTrigger.ActionMap]): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFlagTriggerCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFlagTriggerCommand): CreateFlagTriggerCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFlagTriggerCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFlagTriggerCommand;
  static deserializeBinaryFromReader(message: CreateFlagTriggerCommand, reader: jspb.BinaryReader): CreateFlagTriggerCommand;
}

export namespace CreateFlagTriggerCommand {
  export type AsObject = {
    featureId: string,
    type: proto_autoops_flag_trigger_pb.FlagTrigger.TypeMap[keyof proto_autoops_flag_trigger_pb.FlagTrigger.TypeMap],
    action: proto_autoops_flag_trigger_pb.FlagTrigger.ActionMap[keyof proto_autoops_flag_trigger_pb.FlagTrigger.ActionMap],
    description: string,
  }
}

export class ResetFlagTriggerCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetFlagTriggerCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ResetFlagTriggerCommand): ResetFlagTriggerCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResetFlagTriggerCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetFlagTriggerCommand;
  static deserializeBinaryFromReader(message: ResetFlagTriggerCommand, reader: jspb.BinaryReader): ResetFlagTriggerCommand;
}

export namespace ResetFlagTriggerCommand {
  export type AsObject = {
    id: string,
  }
}

export class ChangeFlagTriggerDescriptionCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeFlagTriggerDescriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeFlagTriggerDescriptionCommand): ChangeFlagTriggerDescriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeFlagTriggerDescriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeFlagTriggerDescriptionCommand;
  static deserializeBinaryFromReader(message: ChangeFlagTriggerDescriptionCommand, reader: jspb.BinaryReader): ChangeFlagTriggerDescriptionCommand;
}

export namespace ChangeFlagTriggerDescriptionCommand {
  export type AsObject = {
    id: string,
    description: string,
  }
}

export class EnableFlagTriggerCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableFlagTriggerCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableFlagTriggerCommand): EnableFlagTriggerCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableFlagTriggerCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableFlagTriggerCommand;
  static deserializeBinaryFromReader(message: EnableFlagTriggerCommand, reader: jspb.BinaryReader): EnableFlagTriggerCommand;
}

export namespace EnableFlagTriggerCommand {
  export type AsObject = {
    id: string,
  }
}

export class DisableFlagTriggerCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableFlagTriggerCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableFlagTriggerCommand): DisableFlagTriggerCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableFlagTriggerCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableFlagTriggerCommand;
  static deserializeBinaryFromReader(message: DisableFlagTriggerCommand, reader: jspb.BinaryReader): DisableFlagTriggerCommand;
}

export namespace DisableFlagTriggerCommand {
  export type AsObject = {
    id: string,
  }
}

export class DeleteFlagTriggerCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteFlagTriggerCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteFlagTriggerCommand): DeleteFlagTriggerCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteFlagTriggerCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteFlagTriggerCommand;
  static deserializeBinaryFromReader(message: DeleteFlagTriggerCommand, reader: jspb.BinaryReader): DeleteFlagTriggerCommand;
}

export namespace DeleteFlagTriggerCommand {
  export type AsObject = {
    id: string,
  }
}

