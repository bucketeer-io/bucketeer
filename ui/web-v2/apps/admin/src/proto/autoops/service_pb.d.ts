// package: bucketeer.autoops
// file: proto/autoops/service.proto

import * as jspb from "google-protobuf";
import * as proto_autoops_auto_ops_rule_pb from "../../proto/autoops/auto_ops_rule_pb";
import * as proto_autoops_command_pb from "../../proto/autoops/command_pb";
import * as proto_autoops_ops_count_pb from "../../proto/autoops/ops_count_pb";
import * as proto_autoops_webhook_pb from "../../proto/autoops/webhook_pb";
import * as proto_autoops_progressive_rollout_pb from "../../proto/autoops/progressive_rollout_pb";

export class GetAutoOpsRuleRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAutoOpsRuleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAutoOpsRuleRequest): GetAutoOpsRuleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAutoOpsRuleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAutoOpsRuleRequest;
  static deserializeBinaryFromReader(message: GetAutoOpsRuleRequest, reader: jspb.BinaryReader): GetAutoOpsRuleRequest;
}

export namespace GetAutoOpsRuleRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
  }
}

export class GetAutoOpsRuleResponse extends jspb.Message {
  hasAutoOpsRule(): boolean;
  clearAutoOpsRule(): void;
  getAutoOpsRule(): proto_autoops_auto_ops_rule_pb.AutoOpsRule | undefined;
  setAutoOpsRule(value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAutoOpsRuleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAutoOpsRuleResponse): GetAutoOpsRuleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAutoOpsRuleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAutoOpsRuleResponse;
  static deserializeBinaryFromReader(message: GetAutoOpsRuleResponse, reader: jspb.BinaryReader): GetAutoOpsRuleResponse;
}

export namespace GetAutoOpsRuleResponse {
  export type AsObject = {
    autoOpsRule?: proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject,
  }
}

export class CreateAutoOpsRuleRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.CreateAutoOpsRuleCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.CreateAutoOpsRuleCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAutoOpsRuleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAutoOpsRuleRequest): CreateAutoOpsRuleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAutoOpsRuleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAutoOpsRuleRequest;
  static deserializeBinaryFromReader(message: CreateAutoOpsRuleRequest, reader: jspb.BinaryReader): CreateAutoOpsRuleRequest;
}

export namespace CreateAutoOpsRuleRequest {
  export type AsObject = {
    environmentNamespace: string,
    command?: proto_autoops_command_pb.CreateAutoOpsRuleCommand.AsObject,
  }
}

export class CreateAutoOpsRuleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAutoOpsRuleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAutoOpsRuleResponse): CreateAutoOpsRuleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAutoOpsRuleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAutoOpsRuleResponse;
  static deserializeBinaryFromReader(message: CreateAutoOpsRuleResponse, reader: jspb.BinaryReader): CreateAutoOpsRuleResponse;
}

export namespace CreateAutoOpsRuleResponse {
  export type AsObject = {
  }
}

export class ListAutoOpsRulesRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearFeatureIdsList(): void;
  getFeatureIdsList(): Array<string>;
  setFeatureIdsList(value: Array<string>): void;
  addFeatureIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAutoOpsRulesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAutoOpsRulesRequest): ListAutoOpsRulesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAutoOpsRulesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAutoOpsRulesRequest;
  static deserializeBinaryFromReader(message: ListAutoOpsRulesRequest, reader: jspb.BinaryReader): ListAutoOpsRulesRequest;
}

export namespace ListAutoOpsRulesRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    featureIdsList: Array<string>,
  }
}

export class ListAutoOpsRulesResponse extends jspb.Message {
  clearAutoOpsRulesList(): void;
  getAutoOpsRulesList(): Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule>;
  setAutoOpsRulesList(value: Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule>): void;
  addAutoOpsRules(value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule, index?: number): proto_autoops_auto_ops_rule_pb.AutoOpsRule;

  getCursor(): string;
  setCursor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAutoOpsRulesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAutoOpsRulesResponse): ListAutoOpsRulesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAutoOpsRulesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAutoOpsRulesResponse;
  static deserializeBinaryFromReader(message: ListAutoOpsRulesResponse, reader: jspb.BinaryReader): ListAutoOpsRulesResponse;
}

export namespace ListAutoOpsRulesResponse {
  export type AsObject = {
    autoOpsRulesList: Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject>,
    cursor: string,
  }
}

export class DeleteAutoOpsRuleRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.DeleteAutoOpsRuleCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.DeleteAutoOpsRuleCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAutoOpsRuleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAutoOpsRuleRequest): DeleteAutoOpsRuleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAutoOpsRuleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAutoOpsRuleRequest;
  static deserializeBinaryFromReader(message: DeleteAutoOpsRuleRequest, reader: jspb.BinaryReader): DeleteAutoOpsRuleRequest;
}

export namespace DeleteAutoOpsRuleRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_autoops_command_pb.DeleteAutoOpsRuleCommand.AsObject,
  }
}

export class DeleteAutoOpsRuleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAutoOpsRuleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAutoOpsRuleResponse): DeleteAutoOpsRuleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAutoOpsRuleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAutoOpsRuleResponse;
  static deserializeBinaryFromReader(message: DeleteAutoOpsRuleResponse, reader: jspb.BinaryReader): DeleteAutoOpsRuleResponse;
}

export namespace DeleteAutoOpsRuleResponse {
  export type AsObject = {
  }
}

export class UpdateAutoOpsRuleRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasChangeAutoOpsRuleOpsTypeCommand(): boolean;
  clearChangeAutoOpsRuleOpsTypeCommand(): void;
  getChangeAutoOpsRuleOpsTypeCommand(): proto_autoops_command_pb.ChangeAutoOpsRuleOpsTypeCommand | undefined;
  setChangeAutoOpsRuleOpsTypeCommand(value?: proto_autoops_command_pb.ChangeAutoOpsRuleOpsTypeCommand): void;

  clearAddOpsEventRateClauseCommandsList(): void;
  getAddOpsEventRateClauseCommandsList(): Array<proto_autoops_command_pb.AddOpsEventRateClauseCommand>;
  setAddOpsEventRateClauseCommandsList(value: Array<proto_autoops_command_pb.AddOpsEventRateClauseCommand>): void;
  addAddOpsEventRateClauseCommands(value?: proto_autoops_command_pb.AddOpsEventRateClauseCommand, index?: number): proto_autoops_command_pb.AddOpsEventRateClauseCommand;

  clearChangeOpsEventRateClauseCommandsList(): void;
  getChangeOpsEventRateClauseCommandsList(): Array<proto_autoops_command_pb.ChangeOpsEventRateClauseCommand>;
  setChangeOpsEventRateClauseCommandsList(value: Array<proto_autoops_command_pb.ChangeOpsEventRateClauseCommand>): void;
  addChangeOpsEventRateClauseCommands(value?: proto_autoops_command_pb.ChangeOpsEventRateClauseCommand, index?: number): proto_autoops_command_pb.ChangeOpsEventRateClauseCommand;

  clearDeleteClauseCommandsList(): void;
  getDeleteClauseCommandsList(): Array<proto_autoops_command_pb.DeleteClauseCommand>;
  setDeleteClauseCommandsList(value: Array<proto_autoops_command_pb.DeleteClauseCommand>): void;
  addDeleteClauseCommands(value?: proto_autoops_command_pb.DeleteClauseCommand, index?: number): proto_autoops_command_pb.DeleteClauseCommand;

  clearAddDatetimeClauseCommandsList(): void;
  getAddDatetimeClauseCommandsList(): Array<proto_autoops_command_pb.AddDatetimeClauseCommand>;
  setAddDatetimeClauseCommandsList(value: Array<proto_autoops_command_pb.AddDatetimeClauseCommand>): void;
  addAddDatetimeClauseCommands(value?: proto_autoops_command_pb.AddDatetimeClauseCommand, index?: number): proto_autoops_command_pb.AddDatetimeClauseCommand;

  clearChangeDatetimeClauseCommandsList(): void;
  getChangeDatetimeClauseCommandsList(): Array<proto_autoops_command_pb.ChangeDatetimeClauseCommand>;
  setChangeDatetimeClauseCommandsList(value: Array<proto_autoops_command_pb.ChangeDatetimeClauseCommand>): void;
  addChangeDatetimeClauseCommands(value?: proto_autoops_command_pb.ChangeDatetimeClauseCommand, index?: number): proto_autoops_command_pb.ChangeDatetimeClauseCommand;

  clearAddWebhookClauseCommandsList(): void;
  getAddWebhookClauseCommandsList(): Array<proto_autoops_command_pb.AddWebhookClauseCommand>;
  setAddWebhookClauseCommandsList(value: Array<proto_autoops_command_pb.AddWebhookClauseCommand>): void;
  addAddWebhookClauseCommands(value?: proto_autoops_command_pb.AddWebhookClauseCommand, index?: number): proto_autoops_command_pb.AddWebhookClauseCommand;

  clearChangeWebhookClauseCommandsList(): void;
  getChangeWebhookClauseCommandsList(): Array<proto_autoops_command_pb.ChangeWebhookClauseCommand>;
  setChangeWebhookClauseCommandsList(value: Array<proto_autoops_command_pb.ChangeWebhookClauseCommand>): void;
  addChangeWebhookClauseCommands(value?: proto_autoops_command_pb.ChangeWebhookClauseCommand, index?: number): proto_autoops_command_pb.ChangeWebhookClauseCommand;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAutoOpsRuleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAutoOpsRuleRequest): UpdateAutoOpsRuleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateAutoOpsRuleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAutoOpsRuleRequest;
  static deserializeBinaryFromReader(message: UpdateAutoOpsRuleRequest, reader: jspb.BinaryReader): UpdateAutoOpsRuleRequest;
}

export namespace UpdateAutoOpsRuleRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    changeAutoOpsRuleOpsTypeCommand?: proto_autoops_command_pb.ChangeAutoOpsRuleOpsTypeCommand.AsObject,
    addOpsEventRateClauseCommandsList: Array<proto_autoops_command_pb.AddOpsEventRateClauseCommand.AsObject>,
    changeOpsEventRateClauseCommandsList: Array<proto_autoops_command_pb.ChangeOpsEventRateClauseCommand.AsObject>,
    deleteClauseCommandsList: Array<proto_autoops_command_pb.DeleteClauseCommand.AsObject>,
    addDatetimeClauseCommandsList: Array<proto_autoops_command_pb.AddDatetimeClauseCommand.AsObject>,
    changeDatetimeClauseCommandsList: Array<proto_autoops_command_pb.ChangeDatetimeClauseCommand.AsObject>,
    addWebhookClauseCommandsList: Array<proto_autoops_command_pb.AddWebhookClauseCommand.AsObject>,
    changeWebhookClauseCommandsList: Array<proto_autoops_command_pb.ChangeWebhookClauseCommand.AsObject>,
  }
}

export class UpdateAutoOpsRuleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAutoOpsRuleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAutoOpsRuleResponse): UpdateAutoOpsRuleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateAutoOpsRuleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAutoOpsRuleResponse;
  static deserializeBinaryFromReader(message: UpdateAutoOpsRuleResponse, reader: jspb.BinaryReader): UpdateAutoOpsRuleResponse;
}

export namespace UpdateAutoOpsRuleResponse {
  export type AsObject = {
  }
}

export class ExecuteAutoOpsRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasChangeAutoOpsRuleTriggeredAtCommand(): boolean;
  clearChangeAutoOpsRuleTriggeredAtCommand(): void;
  getChangeAutoOpsRuleTriggeredAtCommand(): proto_autoops_command_pb.ChangeAutoOpsRuleTriggeredAtCommand | undefined;
  setChangeAutoOpsRuleTriggeredAtCommand(value?: proto_autoops_command_pb.ChangeAutoOpsRuleTriggeredAtCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteAutoOpsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExecuteAutoOpsRequest): ExecuteAutoOpsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecuteAutoOpsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteAutoOpsRequest;
  static deserializeBinaryFromReader(message: ExecuteAutoOpsRequest, reader: jspb.BinaryReader): ExecuteAutoOpsRequest;
}

export namespace ExecuteAutoOpsRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    changeAutoOpsRuleTriggeredAtCommand?: proto_autoops_command_pb.ChangeAutoOpsRuleTriggeredAtCommand.AsObject,
  }
}

export class ExecuteAutoOpsResponse extends jspb.Message {
  getAlreadyTriggered(): boolean;
  setAlreadyTriggered(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteAutoOpsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExecuteAutoOpsResponse): ExecuteAutoOpsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecuteAutoOpsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteAutoOpsResponse;
  static deserializeBinaryFromReader(message: ExecuteAutoOpsResponse, reader: jspb.BinaryReader): ExecuteAutoOpsResponse;
}

export namespace ExecuteAutoOpsResponse {
  export type AsObject = {
    alreadyTriggered: boolean,
  }
}

export class ListOpsCountsRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearAutoOpsRuleIdsList(): void;
  getAutoOpsRuleIdsList(): Array<string>;
  setAutoOpsRuleIdsList(value: Array<string>): void;
  addAutoOpsRuleIds(value: string, index?: number): string;

  clearFeatureIdsList(): void;
  getFeatureIdsList(): Array<string>;
  setFeatureIdsList(value: Array<string>): void;
  addFeatureIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOpsCountsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOpsCountsRequest): ListOpsCountsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListOpsCountsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOpsCountsRequest;
  static deserializeBinaryFromReader(message: ListOpsCountsRequest, reader: jspb.BinaryReader): ListOpsCountsRequest;
}

export namespace ListOpsCountsRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    autoOpsRuleIdsList: Array<string>,
    featureIdsList: Array<string>,
  }
}

export class ListOpsCountsResponse extends jspb.Message {
  getCursor(): string;
  setCursor(value: string): void;

  clearOpsCountsList(): void;
  getOpsCountsList(): Array<proto_autoops_ops_count_pb.OpsCount>;
  setOpsCountsList(value: Array<proto_autoops_ops_count_pb.OpsCount>): void;
  addOpsCounts(value?: proto_autoops_ops_count_pb.OpsCount, index?: number): proto_autoops_ops_count_pb.OpsCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOpsCountsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOpsCountsResponse): ListOpsCountsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListOpsCountsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOpsCountsResponse;
  static deserializeBinaryFromReader(message: ListOpsCountsResponse, reader: jspb.BinaryReader): ListOpsCountsResponse;
}

export namespace ListOpsCountsResponse {
  export type AsObject = {
    cursor: string,
    opsCountsList: Array<proto_autoops_ops_count_pb.OpsCount.AsObject>,
  }
}

export class CreateWebhookRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.CreateWebhookCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.CreateWebhookCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWebhookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWebhookRequest): CreateWebhookRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateWebhookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWebhookRequest;
  static deserializeBinaryFromReader(message: CreateWebhookRequest, reader: jspb.BinaryReader): CreateWebhookRequest;
}

export namespace CreateWebhookRequest {
  export type AsObject = {
    environmentNamespace: string,
    command?: proto_autoops_command_pb.CreateWebhookCommand.AsObject,
  }
}

export class CreateWebhookResponse extends jspb.Message {
  hasWebhook(): boolean;
  clearWebhook(): void;
  getWebhook(): proto_autoops_webhook_pb.Webhook | undefined;
  setWebhook(value?: proto_autoops_webhook_pb.Webhook): void;

  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWebhookResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWebhookResponse): CreateWebhookResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateWebhookResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWebhookResponse;
  static deserializeBinaryFromReader(message: CreateWebhookResponse, reader: jspb.BinaryReader): CreateWebhookResponse;
}

export namespace CreateWebhookResponse {
  export type AsObject = {
    webhook?: proto_autoops_webhook_pb.Webhook.AsObject,
    url: string,
  }
}

export class GetWebhookRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetWebhookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetWebhookRequest): GetWebhookRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetWebhookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetWebhookRequest;
  static deserializeBinaryFromReader(message: GetWebhookRequest, reader: jspb.BinaryReader): GetWebhookRequest;
}

export namespace GetWebhookRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
  }
}

export class GetWebhookResponse extends jspb.Message {
  hasWebhook(): boolean;
  clearWebhook(): void;
  getWebhook(): proto_autoops_webhook_pb.Webhook | undefined;
  setWebhook(value?: proto_autoops_webhook_pb.Webhook): void;

  getUrl(): string;
  setUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetWebhookResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetWebhookResponse): GetWebhookResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetWebhookResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetWebhookResponse;
  static deserializeBinaryFromReader(message: GetWebhookResponse, reader: jspb.BinaryReader): GetWebhookResponse;
}

export namespace GetWebhookResponse {
  export type AsObject = {
    webhook?: proto_autoops_webhook_pb.Webhook.AsObject,
    url: string,
  }
}

export class UpdateWebhookRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasChangewebhooknamecommand(): boolean;
  clearChangewebhooknamecommand(): void;
  getChangewebhooknamecommand(): proto_autoops_command_pb.ChangeWebhookNameCommand | undefined;
  setChangewebhooknamecommand(value?: proto_autoops_command_pb.ChangeWebhookNameCommand): void;

  hasChangewebhookdescriptioncommand(): boolean;
  clearChangewebhookdescriptioncommand(): void;
  getChangewebhookdescriptioncommand(): proto_autoops_command_pb.ChangeWebhookDescriptionCommand | undefined;
  setChangewebhookdescriptioncommand(value?: proto_autoops_command_pb.ChangeWebhookDescriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateWebhookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateWebhookRequest): UpdateWebhookRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateWebhookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateWebhookRequest;
  static deserializeBinaryFromReader(message: UpdateWebhookRequest, reader: jspb.BinaryReader): UpdateWebhookRequest;
}

export namespace UpdateWebhookRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    changewebhooknamecommand?: proto_autoops_command_pb.ChangeWebhookNameCommand.AsObject,
    changewebhookdescriptioncommand?: proto_autoops_command_pb.ChangeWebhookDescriptionCommand.AsObject,
  }
}

export class UpdateWebhookResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateWebhookResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateWebhookResponse): UpdateWebhookResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateWebhookResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateWebhookResponse;
  static deserializeBinaryFromReader(message: UpdateWebhookResponse, reader: jspb.BinaryReader): UpdateWebhookResponse;
}

export namespace UpdateWebhookResponse {
  export type AsObject = {
  }
}

export class DeleteWebhookRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.DeleteWebhookCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.DeleteWebhookCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteWebhookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteWebhookRequest): DeleteWebhookRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteWebhookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteWebhookRequest;
  static deserializeBinaryFromReader(message: DeleteWebhookRequest, reader: jspb.BinaryReader): DeleteWebhookRequest;
}

export namespace DeleteWebhookRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    command?: proto_autoops_command_pb.DeleteWebhookCommand.AsObject,
  }
}

export class DeleteWebhookResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteWebhookResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteWebhookResponse): DeleteWebhookResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteWebhookResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteWebhookResponse;
  static deserializeBinaryFromReader(message: DeleteWebhookResponse, reader: jspb.BinaryReader): DeleteWebhookResponse;
}

export namespace DeleteWebhookResponse {
  export type AsObject = {
  }
}

export class ListWebhooksRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getOrderBy(): ListWebhooksRequest.OrderByMap[keyof ListWebhooksRequest.OrderByMap];
  setOrderBy(value: ListWebhooksRequest.OrderByMap[keyof ListWebhooksRequest.OrderByMap]): void;

  getOrderDirection(): ListWebhooksRequest.OrderDirectionMap[keyof ListWebhooksRequest.OrderDirectionMap];
  setOrderDirection(value: ListWebhooksRequest.OrderDirectionMap[keyof ListWebhooksRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWebhooksRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListWebhooksRequest): ListWebhooksRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListWebhooksRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWebhooksRequest;
  static deserializeBinaryFromReader(message: ListWebhooksRequest, reader: jspb.BinaryReader): ListWebhooksRequest;
}

export namespace ListWebhooksRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    orderBy: ListWebhooksRequest.OrderByMap[keyof ListWebhooksRequest.OrderByMap],
    orderDirection: ListWebhooksRequest.OrderDirectionMap[keyof ListWebhooksRequest.OrderDirectionMap],
    searchKeyword: string,
  }

  export interface OrderByMap {
    DEFAULT: 0;
    NAME: 1;
    CREATED_AT: 2;
    UPDATED_AT: 3;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListWebhooksResponse extends jspb.Message {
  clearWebhooksList(): void;
  getWebhooksList(): Array<proto_autoops_webhook_pb.Webhook>;
  setWebhooksList(value: Array<proto_autoops_webhook_pb.Webhook>): void;
  addWebhooks(value?: proto_autoops_webhook_pb.Webhook, index?: number): proto_autoops_webhook_pb.Webhook;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWebhooksResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListWebhooksResponse): ListWebhooksResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListWebhooksResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWebhooksResponse;
  static deserializeBinaryFromReader(message: ListWebhooksResponse, reader: jspb.BinaryReader): ListWebhooksResponse;
}

export namespace ListWebhooksResponse {
  export type AsObject = {
    webhooksList: Array<proto_autoops_webhook_pb.Webhook.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateProgressiveRolloutRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.CreateProgressiveRolloutCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.CreateProgressiveRolloutCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProgressiveRolloutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateProgressiveRolloutRequest): CreateProgressiveRolloutRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateProgressiveRolloutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateProgressiveRolloutRequest;
  static deserializeBinaryFromReader(message: CreateProgressiveRolloutRequest, reader: jspb.BinaryReader): CreateProgressiveRolloutRequest;
}

export namespace CreateProgressiveRolloutRequest {
  export type AsObject = {
    environmentNamespace: string,
    command?: proto_autoops_command_pb.CreateProgressiveRolloutCommand.AsObject,
  }
}

export class CreateProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProgressiveRolloutResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateProgressiveRolloutResponse): CreateProgressiveRolloutResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateProgressiveRolloutResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateProgressiveRolloutResponse;
  static deserializeBinaryFromReader(message: CreateProgressiveRolloutResponse, reader: jspb.BinaryReader): CreateProgressiveRolloutResponse;
}

export namespace CreateProgressiveRolloutResponse {
  export type AsObject = {
  }
}

export class GetProgressiveRolloutRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProgressiveRolloutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProgressiveRolloutRequest): GetProgressiveRolloutRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProgressiveRolloutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProgressiveRolloutRequest;
  static deserializeBinaryFromReader(message: GetProgressiveRolloutRequest, reader: jspb.BinaryReader): GetProgressiveRolloutRequest;
}

export namespace GetProgressiveRolloutRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
  }
}

export class GetProgressiveRolloutResponse extends jspb.Message {
  hasProgressiveRollout(): boolean;
  clearProgressiveRollout(): void;
  getProgressiveRollout(): proto_autoops_progressive_rollout_pb.ProgressiveRollout | undefined;
  setProgressiveRollout(value?: proto_autoops_progressive_rollout_pb.ProgressiveRollout): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProgressiveRolloutResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProgressiveRolloutResponse): GetProgressiveRolloutResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProgressiveRolloutResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProgressiveRolloutResponse;
  static deserializeBinaryFromReader(message: GetProgressiveRolloutResponse, reader: jspb.BinaryReader): GetProgressiveRolloutResponse;
}

export namespace GetProgressiveRolloutResponse {
  export type AsObject = {
    progressiveRollout?: proto_autoops_progressive_rollout_pb.ProgressiveRollout.AsObject,
  }
}

export class DeleteProgressiveRolloutRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.DeleteProgressiveRolloutCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.DeleteProgressiveRolloutCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProgressiveRolloutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProgressiveRolloutRequest): DeleteProgressiveRolloutRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteProgressiveRolloutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProgressiveRolloutRequest;
  static deserializeBinaryFromReader(message: DeleteProgressiveRolloutRequest, reader: jspb.BinaryReader): DeleteProgressiveRolloutRequest;
}

export namespace DeleteProgressiveRolloutRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_autoops_command_pb.DeleteProgressiveRolloutCommand.AsObject,
  }
}

export class DeleteProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProgressiveRolloutResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProgressiveRolloutResponse): DeleteProgressiveRolloutResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteProgressiveRolloutResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProgressiveRolloutResponse;
  static deserializeBinaryFromReader(message: DeleteProgressiveRolloutResponse, reader: jspb.BinaryReader): DeleteProgressiveRolloutResponse;
}

export namespace DeleteProgressiveRolloutResponse {
  export type AsObject = {
  }
}

export class ListProgressiveRolloutsRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearFeatureIdsList(): void;
  getFeatureIdsList(): Array<string>;
  setFeatureIdsList(value: Array<string>): void;
  addFeatureIds(value: string, index?: number): string;

  getOrderBy(): ListProgressiveRolloutsRequest.OrderByMap[keyof ListProgressiveRolloutsRequest.OrderByMap];
  setOrderBy(value: ListProgressiveRolloutsRequest.OrderByMap[keyof ListProgressiveRolloutsRequest.OrderByMap]): void;

  getOrderDirection(): ListProgressiveRolloutsRequest.OrderDirectionMap[keyof ListProgressiveRolloutsRequest.OrderDirectionMap];
  setOrderDirection(value: ListProgressiveRolloutsRequest.OrderDirectionMap[keyof ListProgressiveRolloutsRequest.OrderDirectionMap]): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap];
  setStatus(value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap]): void;

  hasType(): boolean;
  clearType(): void;
  getType(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap];
  setType(value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProgressiveRolloutsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProgressiveRolloutsRequest): ListProgressiveRolloutsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProgressiveRolloutsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProgressiveRolloutsRequest;
  static deserializeBinaryFromReader(message: ListProgressiveRolloutsRequest, reader: jspb.BinaryReader): ListProgressiveRolloutsRequest;
}

export namespace ListProgressiveRolloutsRequest {
  export type AsObject = {
    environmentNamespace: string,
    pageSize: number,
    cursor: string,
    featureIdsList: Array<string>,
    orderBy: ListProgressiveRolloutsRequest.OrderByMap[keyof ListProgressiveRolloutsRequest.OrderByMap],
    orderDirection: ListProgressiveRolloutsRequest.OrderDirectionMap[keyof ListProgressiveRolloutsRequest.OrderDirectionMap],
    status: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap],
    type: proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap],
  }

  export interface OrderByMap {
    DEFAULT: 0;
    CREATED_AT: 1;
    UPDATED_AT: 2;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListProgressiveRolloutsResponse extends jspb.Message {
  clearProgressiveRolloutsList(): void;
  getProgressiveRolloutsList(): Array<proto_autoops_progressive_rollout_pb.ProgressiveRollout>;
  setProgressiveRolloutsList(value: Array<proto_autoops_progressive_rollout_pb.ProgressiveRollout>): void;
  addProgressiveRollouts(value?: proto_autoops_progressive_rollout_pb.ProgressiveRollout, index?: number): proto_autoops_progressive_rollout_pb.ProgressiveRollout;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProgressiveRolloutsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProgressiveRolloutsResponse): ListProgressiveRolloutsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProgressiveRolloutsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProgressiveRolloutsResponse;
  static deserializeBinaryFromReader(message: ListProgressiveRolloutsResponse, reader: jspb.BinaryReader): ListProgressiveRolloutsResponse;
}

export namespace ListProgressiveRolloutsResponse {
  export type AsObject = {
    progressiveRolloutsList: Array<proto_autoops_progressive_rollout_pb.ProgressiveRollout.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class ExecuteProgressiveRolloutRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasChangeProgressiveRolloutTriggeredAtCommand(): boolean;
  clearChangeProgressiveRolloutTriggeredAtCommand(): void;
  getChangeProgressiveRolloutTriggeredAtCommand(): proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand | undefined;
  setChangeProgressiveRolloutTriggeredAtCommand(value?: proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteProgressiveRolloutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExecuteProgressiveRolloutRequest): ExecuteProgressiveRolloutRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecuteProgressiveRolloutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteProgressiveRolloutRequest;
  static deserializeBinaryFromReader(message: ExecuteProgressiveRolloutRequest, reader: jspb.BinaryReader): ExecuteProgressiveRolloutRequest;
}

export namespace ExecuteProgressiveRolloutRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    changeProgressiveRolloutTriggeredAtCommand?: proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand.AsObject,
  }
}

export class ExecuteProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteProgressiveRolloutResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExecuteProgressiveRolloutResponse): ExecuteProgressiveRolloutResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecuteProgressiveRolloutResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteProgressiveRolloutResponse;
  static deserializeBinaryFromReader(message: ExecuteProgressiveRolloutResponse, reader: jspb.BinaryReader): ExecuteProgressiveRolloutResponse;
}

export namespace ExecuteProgressiveRolloutResponse {
  export type AsObject = {
  }
}

