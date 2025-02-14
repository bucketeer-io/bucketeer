// package: bucketeer.autoops
// file: proto/autoops/service.proto

import * as jspb from 'google-protobuf';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as google_api_field_behavior_pb from '../../google/api/field_behavior_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as proto_autoops_auto_ops_rule_pb from '../../proto/autoops/auto_ops_rule_pb';
import * as proto_autoops_clause_pb from '../../proto/autoops/clause_pb';
import * as proto_autoops_command_pb from '../../proto/autoops/command_pb';
import * as proto_autoops_ops_count_pb from '../../proto/autoops/ops_count_pb';
import * as proto_autoops_progressive_rollout_pb from '../../proto/autoops/progressive_rollout_pb';

export class GetAutoOpsRuleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAutoOpsRuleRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAutoOpsRuleRequest
  ): GetAutoOpsRuleRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAutoOpsRuleRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAutoOpsRuleRequest;
  static deserializeBinaryFromReader(
    message: GetAutoOpsRuleRequest,
    reader: jspb.BinaryReader
  ): GetAutoOpsRuleRequest;
}

export namespace GetAutoOpsRuleRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}

export class GetAutoOpsRuleResponse extends jspb.Message {
  hasAutoOpsRule(): boolean;
  clearAutoOpsRule(): void;
  getAutoOpsRule(): proto_autoops_auto_ops_rule_pb.AutoOpsRule | undefined;
  setAutoOpsRule(value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAutoOpsRuleResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetAutoOpsRuleResponse
  ): GetAutoOpsRuleResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetAutoOpsRuleResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetAutoOpsRuleResponse;
  static deserializeBinaryFromReader(
    message: GetAutoOpsRuleResponse,
    reader: jspb.BinaryReader
  ): GetAutoOpsRuleResponse;
}

export namespace GetAutoOpsRuleResponse {
  export type AsObject = {
    autoOpsRule?: proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject;
  };
}

export class CreateAutoOpsRuleRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.CreateAutoOpsRuleCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.CreateAutoOpsRuleCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

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
  toObject(includeInstance?: boolean): CreateAutoOpsRuleRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAutoOpsRuleRequest
  ): CreateAutoOpsRuleRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAutoOpsRuleRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAutoOpsRuleRequest;
  static deserializeBinaryFromReader(
    message: CreateAutoOpsRuleRequest,
    reader: jspb.BinaryReader
  ): CreateAutoOpsRuleRequest;
}

export namespace CreateAutoOpsRuleRequest {
  export type AsObject = {
    command?: proto_autoops_command_pb.CreateAutoOpsRuleCommand.AsObject;
    environmentId: string;
    featureId: string;
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
    opsEventRateClausesList: Array<proto_autoops_clause_pb.OpsEventRateClause.AsObject>;
    datetimeClausesList: Array<proto_autoops_clause_pb.DatetimeClause.AsObject>;
  };
}

export class CreateAutoOpsRuleResponse extends jspb.Message {
  hasAutoOpsRule(): boolean;
  clearAutoOpsRule(): void;
  getAutoOpsRule(): proto_autoops_auto_ops_rule_pb.AutoOpsRule | undefined;
  setAutoOpsRule(value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAutoOpsRuleResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateAutoOpsRuleResponse
  ): CreateAutoOpsRuleResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateAutoOpsRuleResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateAutoOpsRuleResponse;
  static deserializeBinaryFromReader(
    message: CreateAutoOpsRuleResponse,
    reader: jspb.BinaryReader
  ): CreateAutoOpsRuleResponse;
}

export namespace CreateAutoOpsRuleResponse {
  export type AsObject = {
    autoOpsRule?: proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject;
  };
}

export class ListAutoOpsRulesRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearFeatureIdsList(): void;
  getFeatureIdsList(): Array<string>;
  setFeatureIdsList(value: Array<string>): void;
  addFeatureIds(value: string, index?: number): string;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAutoOpsRulesRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListAutoOpsRulesRequest
  ): ListAutoOpsRulesRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListAutoOpsRulesRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListAutoOpsRulesRequest;
  static deserializeBinaryFromReader(
    message: ListAutoOpsRulesRequest,
    reader: jspb.BinaryReader
  ): ListAutoOpsRulesRequest;
}

export namespace ListAutoOpsRulesRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    featureIdsList: Array<string>;
    environmentId: string;
  };
}

export class ListAutoOpsRulesResponse extends jspb.Message {
  clearAutoOpsRulesList(): void;
  getAutoOpsRulesList(): Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule>;
  setAutoOpsRulesList(
    value: Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule>
  ): void;
  addAutoOpsRules(
    value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule,
    index?: number
  ): proto_autoops_auto_ops_rule_pb.AutoOpsRule;

  getCursor(): string;
  setCursor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAutoOpsRulesResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListAutoOpsRulesResponse
  ): ListAutoOpsRulesResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListAutoOpsRulesResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListAutoOpsRulesResponse;
  static deserializeBinaryFromReader(
    message: ListAutoOpsRulesResponse,
    reader: jspb.BinaryReader
  ): ListAutoOpsRulesResponse;
}

export namespace ListAutoOpsRulesResponse {
  export type AsObject = {
    autoOpsRulesList: Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject>;
    cursor: string;
  };
}

export class StopAutoOpsRuleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.StopAutoOpsRuleCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.StopAutoOpsRuleCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopAutoOpsRuleRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StopAutoOpsRuleRequest
  ): StopAutoOpsRuleRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StopAutoOpsRuleRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StopAutoOpsRuleRequest;
  static deserializeBinaryFromReader(
    message: StopAutoOpsRuleRequest,
    reader: jspb.BinaryReader
  ): StopAutoOpsRuleRequest;
}

export namespace StopAutoOpsRuleRequest {
  export type AsObject = {
    id: string;
    command?: proto_autoops_command_pb.StopAutoOpsRuleCommand.AsObject;
    environmentId: string;
  };
}

export class StopAutoOpsRuleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopAutoOpsRuleResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StopAutoOpsRuleResponse
  ): StopAutoOpsRuleResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StopAutoOpsRuleResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StopAutoOpsRuleResponse;
  static deserializeBinaryFromReader(
    message: StopAutoOpsRuleResponse,
    reader: jspb.BinaryReader
  ): StopAutoOpsRuleResponse;
}

export namespace StopAutoOpsRuleResponse {
  export type AsObject = {};
}

export class DeleteAutoOpsRuleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_autoops_command_pb.DeleteAutoOpsRuleCommand | undefined;
  setCommand(value?: proto_autoops_command_pb.DeleteAutoOpsRuleCommand): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAutoOpsRuleRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteAutoOpsRuleRequest
  ): DeleteAutoOpsRuleRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteAutoOpsRuleRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAutoOpsRuleRequest;
  static deserializeBinaryFromReader(
    message: DeleteAutoOpsRuleRequest,
    reader: jspb.BinaryReader
  ): DeleteAutoOpsRuleRequest;
}

export namespace DeleteAutoOpsRuleRequest {
  export type AsObject = {
    id: string;
    command?: proto_autoops_command_pb.DeleteAutoOpsRuleCommand.AsObject;
    environmentId: string;
  };
}

export class DeleteAutoOpsRuleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAutoOpsRuleResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteAutoOpsRuleResponse
  ): DeleteAutoOpsRuleResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteAutoOpsRuleResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAutoOpsRuleResponse;
  static deserializeBinaryFromReader(
    message: DeleteAutoOpsRuleResponse,
    reader: jspb.BinaryReader
  ): DeleteAutoOpsRuleResponse;
}

export namespace DeleteAutoOpsRuleResponse {
  export type AsObject = {};
}

export class UpdateAutoOpsRuleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  clearAddOpsEventRateClauseCommandsList(): void;
  getAddOpsEventRateClauseCommandsList(): Array<proto_autoops_command_pb.AddOpsEventRateClauseCommand>;
  setAddOpsEventRateClauseCommandsList(
    value: Array<proto_autoops_command_pb.AddOpsEventRateClauseCommand>
  ): void;
  addAddOpsEventRateClauseCommands(
    value?: proto_autoops_command_pb.AddOpsEventRateClauseCommand,
    index?: number
  ): proto_autoops_command_pb.AddOpsEventRateClauseCommand;

  clearChangeOpsEventRateClauseCommandsList(): void;
  getChangeOpsEventRateClauseCommandsList(): Array<proto_autoops_command_pb.ChangeOpsEventRateClauseCommand>;
  setChangeOpsEventRateClauseCommandsList(
    value: Array<proto_autoops_command_pb.ChangeOpsEventRateClauseCommand>
  ): void;
  addChangeOpsEventRateClauseCommands(
    value?: proto_autoops_command_pb.ChangeOpsEventRateClauseCommand,
    index?: number
  ): proto_autoops_command_pb.ChangeOpsEventRateClauseCommand;

  clearDeleteClauseCommandsList(): void;
  getDeleteClauseCommandsList(): Array<proto_autoops_command_pb.DeleteClauseCommand>;
  setDeleteClauseCommandsList(
    value: Array<proto_autoops_command_pb.DeleteClauseCommand>
  ): void;
  addDeleteClauseCommands(
    value?: proto_autoops_command_pb.DeleteClauseCommand,
    index?: number
  ): proto_autoops_command_pb.DeleteClauseCommand;

  clearAddDatetimeClauseCommandsList(): void;
  getAddDatetimeClauseCommandsList(): Array<proto_autoops_command_pb.AddDatetimeClauseCommand>;
  setAddDatetimeClauseCommandsList(
    value: Array<proto_autoops_command_pb.AddDatetimeClauseCommand>
  ): void;
  addAddDatetimeClauseCommands(
    value?: proto_autoops_command_pb.AddDatetimeClauseCommand,
    index?: number
  ): proto_autoops_command_pb.AddDatetimeClauseCommand;

  clearChangeDatetimeClauseCommandsList(): void;
  getChangeDatetimeClauseCommandsList(): Array<proto_autoops_command_pb.ChangeDatetimeClauseCommand>;
  setChangeDatetimeClauseCommandsList(
    value: Array<proto_autoops_command_pb.ChangeDatetimeClauseCommand>
  ): void;
  addChangeDatetimeClauseCommands(
    value?: proto_autoops_command_pb.ChangeDatetimeClauseCommand,
    index?: number
  ): proto_autoops_command_pb.ChangeDatetimeClauseCommand;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  clearUpdateOpsEventRateClausesList(): void;
  getUpdateOpsEventRateClausesList(): Array<UpdateAutoOpsRuleRequest.UpdateOpsEventRateClause>;
  setUpdateOpsEventRateClausesList(
    value: Array<UpdateAutoOpsRuleRequest.UpdateOpsEventRateClause>
  ): void;
  addUpdateOpsEventRateClauses(
    value?: UpdateAutoOpsRuleRequest.UpdateOpsEventRateClause,
    index?: number
  ): UpdateAutoOpsRuleRequest.UpdateOpsEventRateClause;

  clearUpdateDatetimeClausesList(): void;
  getUpdateDatetimeClausesList(): Array<UpdateAutoOpsRuleRequest.UpdateDatetimeClause>;
  setUpdateDatetimeClausesList(
    value: Array<UpdateAutoOpsRuleRequest.UpdateDatetimeClause>
  ): void;
  addUpdateDatetimeClauses(
    value?: UpdateAutoOpsRuleRequest.UpdateDatetimeClause,
    index?: number
  ): UpdateAutoOpsRuleRequest.UpdateDatetimeClause;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAutoOpsRuleRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateAutoOpsRuleRequest
  ): UpdateAutoOpsRuleRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateAutoOpsRuleRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAutoOpsRuleRequest;
  static deserializeBinaryFromReader(
    message: UpdateAutoOpsRuleRequest,
    reader: jspb.BinaryReader
  ): UpdateAutoOpsRuleRequest;
}

export namespace UpdateAutoOpsRuleRequest {
  export type AsObject = {
    id: string;
    addOpsEventRateClauseCommandsList: Array<proto_autoops_command_pb.AddOpsEventRateClauseCommand.AsObject>;
    changeOpsEventRateClauseCommandsList: Array<proto_autoops_command_pb.ChangeOpsEventRateClauseCommand.AsObject>;
    deleteClauseCommandsList: Array<proto_autoops_command_pb.DeleteClauseCommand.AsObject>;
    addDatetimeClauseCommandsList: Array<proto_autoops_command_pb.AddDatetimeClauseCommand.AsObject>;
    changeDatetimeClauseCommandsList: Array<proto_autoops_command_pb.ChangeDatetimeClauseCommand.AsObject>;
    environmentId: string;
    updateOpsEventRateClausesList: Array<UpdateAutoOpsRuleRequest.UpdateOpsEventRateClause.AsObject>;
    updateDatetimeClausesList: Array<UpdateAutoOpsRuleRequest.UpdateDatetimeClause.AsObject>;
  };

  export class UpdateOpsEventRateClause extends jspb.Message {
    getId(): string;
    setId(value: string): void;

    hasDeleted(): boolean;
    clearDeleted(): void;
    getDeleted(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setDeleted(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasClause(): boolean;
    clearClause(): void;
    getClause(): proto_autoops_clause_pb.OpsEventRateClause | undefined;
    setClause(value?: proto_autoops_clause_pb.OpsEventRateClause): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateOpsEventRateClause.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: UpdateOpsEventRateClause
    ): UpdateOpsEventRateClause.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: UpdateOpsEventRateClause,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): UpdateOpsEventRateClause;
    static deserializeBinaryFromReader(
      message: UpdateOpsEventRateClause,
      reader: jspb.BinaryReader
    ): UpdateOpsEventRateClause;
  }

  export namespace UpdateOpsEventRateClause {
    export type AsObject = {
      id: string;
      deleted?: google_protobuf_wrappers_pb.BoolValue.AsObject;
      clause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject;
    };
  }

  export class UpdateDatetimeClause extends jspb.Message {
    getId(): string;
    setId(value: string): void;

    hasDeleted(): boolean;
    clearDeleted(): void;
    getDeleted(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setDeleted(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasClause(): boolean;
    clearClause(): void;
    getClause(): proto_autoops_clause_pb.DatetimeClause | undefined;
    setClause(value?: proto_autoops_clause_pb.DatetimeClause): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateDatetimeClause.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: UpdateDatetimeClause
    ): UpdateDatetimeClause.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: UpdateDatetimeClause,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): UpdateDatetimeClause;
    static deserializeBinaryFromReader(
      message: UpdateDatetimeClause,
      reader: jspb.BinaryReader
    ): UpdateDatetimeClause;
  }

  export namespace UpdateDatetimeClause {
    export type AsObject = {
      id: string;
      deleted?: google_protobuf_wrappers_pb.BoolValue.AsObject;
      clause?: proto_autoops_clause_pb.DatetimeClause.AsObject;
    };
  }
}

export class UpdateAutoOpsRuleResponse extends jspb.Message {
  hasAutoOpsRule(): boolean;
  clearAutoOpsRule(): void;
  getAutoOpsRule(): proto_autoops_auto_ops_rule_pb.AutoOpsRule | undefined;
  setAutoOpsRule(value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAutoOpsRuleResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateAutoOpsRuleResponse
  ): UpdateAutoOpsRuleResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateAutoOpsRuleResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAutoOpsRuleResponse;
  static deserializeBinaryFromReader(
    message: UpdateAutoOpsRuleResponse,
    reader: jspb.BinaryReader
  ): UpdateAutoOpsRuleResponse;
}

export namespace UpdateAutoOpsRuleResponse {
  export type AsObject = {
    autoOpsRule?: proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject;
  };
}

export class ExecuteAutoOpsRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasExecuteAutoOpsRuleCommand(): boolean;
  clearExecuteAutoOpsRuleCommand(): void;
  getExecuteAutoOpsRuleCommand():
    | proto_autoops_command_pb.ExecuteAutoOpsRuleCommand
    | undefined;
  setExecuteAutoOpsRuleCommand(
    value?: proto_autoops_command_pb.ExecuteAutoOpsRuleCommand
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteAutoOpsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExecuteAutoOpsRequest
  ): ExecuteAutoOpsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExecuteAutoOpsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteAutoOpsRequest;
  static deserializeBinaryFromReader(
    message: ExecuteAutoOpsRequest,
    reader: jspb.BinaryReader
  ): ExecuteAutoOpsRequest;
}

export namespace ExecuteAutoOpsRequest {
  export type AsObject = {
    id: string;
    executeAutoOpsRuleCommand?: proto_autoops_command_pb.ExecuteAutoOpsRuleCommand.AsObject;
    environmentId: string;
  };
}

export class ExecuteAutoOpsResponse extends jspb.Message {
  getAlreadyTriggered(): boolean;
  setAlreadyTriggered(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteAutoOpsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExecuteAutoOpsResponse
  ): ExecuteAutoOpsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExecuteAutoOpsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteAutoOpsResponse;
  static deserializeBinaryFromReader(
    message: ExecuteAutoOpsResponse,
    reader: jspb.BinaryReader
  ): ExecuteAutoOpsResponse;
}

export namespace ExecuteAutoOpsResponse {
  export type AsObject = {
    alreadyTriggered: boolean;
  };
}

export class ListOpsCountsRequest extends jspb.Message {
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

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOpsCountsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListOpsCountsRequest
  ): ListOpsCountsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListOpsCountsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListOpsCountsRequest;
  static deserializeBinaryFromReader(
    message: ListOpsCountsRequest,
    reader: jspb.BinaryReader
  ): ListOpsCountsRequest;
}

export namespace ListOpsCountsRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    autoOpsRuleIdsList: Array<string>;
    featureIdsList: Array<string>;
    environmentId: string;
  };
}

export class ListOpsCountsResponse extends jspb.Message {
  getCursor(): string;
  setCursor(value: string): void;

  clearOpsCountsList(): void;
  getOpsCountsList(): Array<proto_autoops_ops_count_pb.OpsCount>;
  setOpsCountsList(value: Array<proto_autoops_ops_count_pb.OpsCount>): void;
  addOpsCounts(
    value?: proto_autoops_ops_count_pb.OpsCount,
    index?: number
  ): proto_autoops_ops_count_pb.OpsCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOpsCountsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListOpsCountsResponse
  ): ListOpsCountsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListOpsCountsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListOpsCountsResponse;
  static deserializeBinaryFromReader(
    message: ListOpsCountsResponse,
    reader: jspb.BinaryReader
  ): ListOpsCountsResponse;
}

export namespace ListOpsCountsResponse {
  export type AsObject = {
    cursor: string;
    opsCountsList: Array<proto_autoops_ops_count_pb.OpsCount.AsObject>;
  };
}

export class CreateProgressiveRolloutRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand():
    | proto_autoops_command_pb.CreateProgressiveRolloutCommand
    | undefined;
  setCommand(
    value?: proto_autoops_command_pb.CreateProgressiveRolloutCommand
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProgressiveRolloutRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateProgressiveRolloutRequest
  ): CreateProgressiveRolloutRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateProgressiveRolloutRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateProgressiveRolloutRequest;
  static deserializeBinaryFromReader(
    message: CreateProgressiveRolloutRequest,
    reader: jspb.BinaryReader
  ): CreateProgressiveRolloutRequest;
}

export namespace CreateProgressiveRolloutRequest {
  export type AsObject = {
    command?: proto_autoops_command_pb.CreateProgressiveRolloutCommand.AsObject;
    environmentId: string;
  };
}

export class CreateProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): CreateProgressiveRolloutResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateProgressiveRolloutResponse
  ): CreateProgressiveRolloutResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateProgressiveRolloutResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateProgressiveRolloutResponse;
  static deserializeBinaryFromReader(
    message: CreateProgressiveRolloutResponse,
    reader: jspb.BinaryReader
  ): CreateProgressiveRolloutResponse;
}

export namespace CreateProgressiveRolloutResponse {
  export type AsObject = {};
}

export class GetProgressiveRolloutRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProgressiveRolloutRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetProgressiveRolloutRequest
  ): GetProgressiveRolloutRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetProgressiveRolloutRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetProgressiveRolloutRequest;
  static deserializeBinaryFromReader(
    message: GetProgressiveRolloutRequest,
    reader: jspb.BinaryReader
  ): GetProgressiveRolloutRequest;
}

export namespace GetProgressiveRolloutRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}

export class GetProgressiveRolloutResponse extends jspb.Message {
  hasProgressiveRollout(): boolean;
  clearProgressiveRollout(): void;
  getProgressiveRollout():
    | proto_autoops_progressive_rollout_pb.ProgressiveRollout
    | undefined;
  setProgressiveRollout(
    value?: proto_autoops_progressive_rollout_pb.ProgressiveRollout
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProgressiveRolloutResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetProgressiveRolloutResponse
  ): GetProgressiveRolloutResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetProgressiveRolloutResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetProgressiveRolloutResponse;
  static deserializeBinaryFromReader(
    message: GetProgressiveRolloutResponse,
    reader: jspb.BinaryReader
  ): GetProgressiveRolloutResponse;
}

export namespace GetProgressiveRolloutResponse {
  export type AsObject = {
    progressiveRollout?: proto_autoops_progressive_rollout_pb.ProgressiveRollout.AsObject;
  };
}

export class StopProgressiveRolloutRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand():
    | proto_autoops_command_pb.StopProgressiveRolloutCommand
    | undefined;
  setCommand(
    value?: proto_autoops_command_pb.StopProgressiveRolloutCommand
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopProgressiveRolloutRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StopProgressiveRolloutRequest
  ): StopProgressiveRolloutRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StopProgressiveRolloutRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StopProgressiveRolloutRequest;
  static deserializeBinaryFromReader(
    message: StopProgressiveRolloutRequest,
    reader: jspb.BinaryReader
  ): StopProgressiveRolloutRequest;
}

export namespace StopProgressiveRolloutRequest {
  export type AsObject = {
    id: string;
    command?: proto_autoops_command_pb.StopProgressiveRolloutCommand.AsObject;
    environmentId: string;
  };
}

export class StopProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopProgressiveRolloutResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StopProgressiveRolloutResponse
  ): StopProgressiveRolloutResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StopProgressiveRolloutResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StopProgressiveRolloutResponse;
  static deserializeBinaryFromReader(
    message: StopProgressiveRolloutResponse,
    reader: jspb.BinaryReader
  ): StopProgressiveRolloutResponse;
}

export namespace StopProgressiveRolloutResponse {
  export type AsObject = {};
}

export class DeleteProgressiveRolloutRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand():
    | proto_autoops_command_pb.DeleteProgressiveRolloutCommand
    | undefined;
  setCommand(
    value?: proto_autoops_command_pb.DeleteProgressiveRolloutCommand
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProgressiveRolloutRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteProgressiveRolloutRequest
  ): DeleteProgressiveRolloutRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteProgressiveRolloutRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProgressiveRolloutRequest;
  static deserializeBinaryFromReader(
    message: DeleteProgressiveRolloutRequest,
    reader: jspb.BinaryReader
  ): DeleteProgressiveRolloutRequest;
}

export namespace DeleteProgressiveRolloutRequest {
  export type AsObject = {
    id: string;
    command?: proto_autoops_command_pb.DeleteProgressiveRolloutCommand.AsObject;
    environmentId: string;
  };
}

export class DeleteProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): DeleteProgressiveRolloutResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteProgressiveRolloutResponse
  ): DeleteProgressiveRolloutResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteProgressiveRolloutResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProgressiveRolloutResponse;
  static deserializeBinaryFromReader(
    message: DeleteProgressiveRolloutResponse,
    reader: jspb.BinaryReader
  ): DeleteProgressiveRolloutResponse;
}

export namespace DeleteProgressiveRolloutResponse {
  export type AsObject = {};
}

export class ListProgressiveRolloutsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  clearFeatureIdsList(): void;
  getFeatureIdsList(): Array<string>;
  setFeatureIdsList(value: Array<string>): void;
  addFeatureIds(value: string, index?: number): string;

  getOrderBy(): ListProgressiveRolloutsRequest.OrderByMap[keyof ListProgressiveRolloutsRequest.OrderByMap];
  setOrderBy(
    value: ListProgressiveRolloutsRequest.OrderByMap[keyof ListProgressiveRolloutsRequest.OrderByMap]
  ): void;

  getOrderDirection(): ListProgressiveRolloutsRequest.OrderDirectionMap[keyof ListProgressiveRolloutsRequest.OrderDirectionMap];
  setOrderDirection(
    value: ListProgressiveRolloutsRequest.OrderDirectionMap[keyof ListProgressiveRolloutsRequest.OrderDirectionMap]
  ): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap];
  setStatus(
    value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap]
  ): void;

  hasType(): boolean;
  clearType(): void;
  getType(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap];
  setType(
    value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap]
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProgressiveRolloutsRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListProgressiveRolloutsRequest
  ): ListProgressiveRolloutsRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListProgressiveRolloutsRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListProgressiveRolloutsRequest;
  static deserializeBinaryFromReader(
    message: ListProgressiveRolloutsRequest,
    reader: jspb.BinaryReader
  ): ListProgressiveRolloutsRequest;
}

export namespace ListProgressiveRolloutsRequest {
  export type AsObject = {
    pageSize: number;
    cursor: string;
    featureIdsList: Array<string>;
    orderBy: ListProgressiveRolloutsRequest.OrderByMap[keyof ListProgressiveRolloutsRequest.OrderByMap];
    orderDirection: ListProgressiveRolloutsRequest.OrderDirectionMap[keyof ListProgressiveRolloutsRequest.OrderDirectionMap];
    status: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap];
    type: proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap];
    environmentId: string;
  };

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
  setProgressiveRolloutsList(
    value: Array<proto_autoops_progressive_rollout_pb.ProgressiveRollout>
  ): void;
  addProgressiveRollouts(
    value?: proto_autoops_progressive_rollout_pb.ProgressiveRollout,
    index?: number
  ): proto_autoops_progressive_rollout_pb.ProgressiveRollout;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProgressiveRolloutsResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListProgressiveRolloutsResponse
  ): ListProgressiveRolloutsResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListProgressiveRolloutsResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListProgressiveRolloutsResponse;
  static deserializeBinaryFromReader(
    message: ListProgressiveRolloutsResponse,
    reader: jspb.BinaryReader
  ): ListProgressiveRolloutsResponse;
}

export namespace ListProgressiveRolloutsResponse {
  export type AsObject = {
    progressiveRolloutsList: Array<proto_autoops_progressive_rollout_pb.ProgressiveRollout.AsObject>;
    cursor: string;
    totalCount: number;
  };
}

export class ExecuteProgressiveRolloutRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasChangeProgressiveRolloutTriggeredAtCommand(): boolean;
  clearChangeProgressiveRolloutTriggeredAtCommand(): void;
  getChangeProgressiveRolloutTriggeredAtCommand():
    | proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand
    | undefined;
  setChangeProgressiveRolloutTriggeredAtCommand(
    value?: proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ExecuteProgressiveRolloutRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExecuteProgressiveRolloutRequest
  ): ExecuteProgressiveRolloutRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExecuteProgressiveRolloutRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteProgressiveRolloutRequest;
  static deserializeBinaryFromReader(
    message: ExecuteProgressiveRolloutRequest,
    reader: jspb.BinaryReader
  ): ExecuteProgressiveRolloutRequest;
}

export namespace ExecuteProgressiveRolloutRequest {
  export type AsObject = {
    id: string;
    changeProgressiveRolloutTriggeredAtCommand?: proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand.AsObject;
    environmentId: string;
  };
}

export class ExecuteProgressiveRolloutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ExecuteProgressiveRolloutResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExecuteProgressiveRolloutResponse
  ): ExecuteProgressiveRolloutResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExecuteProgressiveRolloutResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ExecuteProgressiveRolloutResponse;
  static deserializeBinaryFromReader(
    message: ExecuteProgressiveRolloutResponse,
    reader: jspb.BinaryReader
  ): ExecuteProgressiveRolloutResponse;
}

export namespace ExecuteProgressiveRolloutResponse {
  export type AsObject = {};
}
