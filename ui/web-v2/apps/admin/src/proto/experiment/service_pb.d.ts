// package: bucketeer.experiment
// file: proto/experiment/service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as proto_experiment_command_pb from "../../proto/experiment/command_pb";
import * as proto_experiment_goal_pb from "../../proto/experiment/goal_pb";
import * as proto_experiment_experiment_pb from "../../proto/experiment/experiment_pb";

export class GetGoalRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGoalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGoalRequest): GetGoalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGoalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGoalRequest;
  static deserializeBinaryFromReader(message: GetGoalRequest, reader: jspb.BinaryReader): GetGoalRequest;
}

export namespace GetGoalRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
  }
}

export class GetGoalResponse extends jspb.Message {
  hasGoal(): boolean;
  clearGoal(): void;
  getGoal(): proto_experiment_goal_pb.Goal | undefined;
  setGoal(value?: proto_experiment_goal_pb.Goal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGoalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGoalResponse): GetGoalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGoalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGoalResponse;
  static deserializeBinaryFromReader(message: GetGoalResponse, reader: jspb.BinaryReader): GetGoalResponse;
}

export namespace GetGoalResponse {
  export type AsObject = {
    goal?: proto_experiment_goal_pb.Goal.AsObject,
  }
}

export class ListGoalsRequest extends jspb.Message {
  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getOrderBy(): ListGoalsRequest.OrderByMap[keyof ListGoalsRequest.OrderByMap];
  setOrderBy(value: ListGoalsRequest.OrderByMap[keyof ListGoalsRequest.OrderByMap]): void;

  getOrderDirection(): ListGoalsRequest.OrderDirectionMap[keyof ListGoalsRequest.OrderDirectionMap];
  setOrderDirection(value: ListGoalsRequest.OrderDirectionMap[keyof ListGoalsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasIsInUseStatus(): boolean;
  clearIsInUseStatus(): void;
  getIsInUseStatus(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIsInUseStatus(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasArchived(): boolean;
  clearArchived(): void;
  getArchived(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setArchived(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGoalsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGoalsRequest): ListGoalsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGoalsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGoalsRequest;
  static deserializeBinaryFromReader(message: ListGoalsRequest, reader: jspb.BinaryReader): ListGoalsRequest;
}

export namespace ListGoalsRequest {
  export type AsObject = {
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    orderBy: ListGoalsRequest.OrderByMap[keyof ListGoalsRequest.OrderByMap],
    orderDirection: ListGoalsRequest.OrderDirectionMap[keyof ListGoalsRequest.OrderDirectionMap],
    searchKeyword: string,
    isInUseStatus?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    archived?: google_protobuf_wrappers_pb.BoolValue.AsObject,
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

export class ListGoalsResponse extends jspb.Message {
  clearGoalsList(): void;
  getGoalsList(): Array<proto_experiment_goal_pb.Goal>;
  setGoalsList(value: Array<proto_experiment_goal_pb.Goal>): void;
  addGoals(value?: proto_experiment_goal_pb.Goal, index?: number): proto_experiment_goal_pb.Goal;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGoalsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGoalsResponse): ListGoalsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGoalsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGoalsResponse;
  static deserializeBinaryFromReader(message: ListGoalsResponse, reader: jspb.BinaryReader): ListGoalsResponse;
}

export namespace ListGoalsResponse {
  export type AsObject = {
    goalsList: Array<proto_experiment_goal_pb.Goal.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateGoalRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.CreateGoalCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.CreateGoalCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGoalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGoalRequest): CreateGoalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGoalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGoalRequest;
  static deserializeBinaryFromReader(message: CreateGoalRequest, reader: jspb.BinaryReader): CreateGoalRequest;
}

export namespace CreateGoalRequest {
  export type AsObject = {
    command?: proto_experiment_command_pb.CreateGoalCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CreateGoalResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGoalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGoalResponse): CreateGoalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGoalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGoalResponse;
  static deserializeBinaryFromReader(message: CreateGoalResponse, reader: jspb.BinaryReader): CreateGoalResponse;
}

export namespace CreateGoalResponse {
  export type AsObject = {
  }
}

export class ArchiveGoalRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.ArchiveGoalCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.ArchiveGoalCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveGoalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveGoalRequest): ArchiveGoalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveGoalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveGoalRequest;
  static deserializeBinaryFromReader(message: ArchiveGoalRequest, reader: jspb.BinaryReader): ArchiveGoalRequest;
}

export namespace ArchiveGoalRequest {
  export type AsObject = {
    id: string,
    command?: proto_experiment_command_pb.ArchiveGoalCommand.AsObject,
    environmentNamespace: string,
  }
}

export class ArchiveGoalResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveGoalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveGoalResponse): ArchiveGoalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveGoalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveGoalResponse;
  static deserializeBinaryFromReader(message: ArchiveGoalResponse, reader: jspb.BinaryReader): ArchiveGoalResponse;
}

export namespace ArchiveGoalResponse {
  export type AsObject = {
  }
}

export class DeleteGoalRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.DeleteGoalCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.DeleteGoalCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGoalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGoalRequest): DeleteGoalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGoalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGoalRequest;
  static deserializeBinaryFromReader(message: DeleteGoalRequest, reader: jspb.BinaryReader): DeleteGoalRequest;
}

export namespace DeleteGoalRequest {
  export type AsObject = {
    id: string,
    command?: proto_experiment_command_pb.DeleteGoalCommand.AsObject,
    environmentNamespace: string,
  }
}

export class DeleteGoalResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGoalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGoalResponse): DeleteGoalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGoalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGoalResponse;
  static deserializeBinaryFromReader(message: DeleteGoalResponse, reader: jspb.BinaryReader): DeleteGoalResponse;
}

export namespace DeleteGoalResponse {
  export type AsObject = {
  }
}

export class UpdateGoalRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRenameCommand(): boolean;
  clearRenameCommand(): void;
  getRenameCommand(): proto_experiment_command_pb.RenameGoalCommand | undefined;
  setRenameCommand(value?: proto_experiment_command_pb.RenameGoalCommand): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_experiment_command_pb.ChangeDescriptionGoalCommand | undefined;
  setChangeDescriptionCommand(value?: proto_experiment_command_pb.ChangeDescriptionGoalCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGoalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGoalRequest): UpdateGoalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGoalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGoalRequest;
  static deserializeBinaryFromReader(message: UpdateGoalRequest, reader: jspb.BinaryReader): UpdateGoalRequest;
}

export namespace UpdateGoalRequest {
  export type AsObject = {
    id: string,
    renameCommand?: proto_experiment_command_pb.RenameGoalCommand.AsObject,
    changeDescriptionCommand?: proto_experiment_command_pb.ChangeDescriptionGoalCommand.AsObject,
    environmentNamespace: string,
  }
}

export class UpdateGoalResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGoalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGoalResponse): UpdateGoalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGoalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGoalResponse;
  static deserializeBinaryFromReader(message: UpdateGoalResponse, reader: jspb.BinaryReader): UpdateGoalResponse;
}

export namespace UpdateGoalResponse {
  export type AsObject = {
  }
}

export class GetExperimentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentRequest): GetExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentRequest;
  static deserializeBinaryFromReader(message: GetExperimentRequest, reader: jspb.BinaryReader): GetExperimentRequest;
}

export namespace GetExperimentRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
  }
}

export class GetExperimentResponse extends jspb.Message {
  hasExperiment(): boolean;
  clearExperiment(): void;
  getExperiment(): proto_experiment_experiment_pb.Experiment | undefined;
  setExperiment(value?: proto_experiment_experiment_pb.Experiment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetExperimentResponse): GetExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetExperimentResponse;
  static deserializeBinaryFromReader(message: GetExperimentResponse, reader: jspb.BinaryReader): GetExperimentResponse;
}

export namespace GetExperimentResponse {
  export type AsObject = {
    experiment?: proto_experiment_experiment_pb.Experiment.AsObject,
  }
}

export class ListExperimentsRequest extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  hasFeatureVersion(): boolean;
  clearFeatureVersion(): void;
  getFeatureVersion(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setFeatureVersion(value?: google_protobuf_wrappers_pb.Int32Value): void;

  getFrom(): number;
  setFrom(value: number): void;

  getTo(): number;
  setTo(value: number): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getCursor(): string;
  setCursor(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setStatus(value?: google_protobuf_wrappers_pb.Int32Value): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getOrderBy(): ListExperimentsRequest.OrderByMap[keyof ListExperimentsRequest.OrderByMap];
  setOrderBy(value: ListExperimentsRequest.OrderByMap[keyof ListExperimentsRequest.OrderByMap]): void;

  getOrderDirection(): ListExperimentsRequest.OrderDirectionMap[keyof ListExperimentsRequest.OrderDirectionMap];
  setOrderDirection(value: ListExperimentsRequest.OrderDirectionMap[keyof ListExperimentsRequest.OrderDirectionMap]): void;

  getSearchKeyword(): string;
  setSearchKeyword(value: string): void;

  hasArchived(): boolean;
  clearArchived(): void;
  getArchived(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setArchived(value?: google_protobuf_wrappers_pb.BoolValue): void;

  clearStatusesList(): void;
  getStatusesList(): Array<proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap]>;
  setStatusesList(value: Array<proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap]>): void;
  addStatuses(value: proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap], index?: number): proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap];

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExperimentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListExperimentsRequest): ListExperimentsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListExperimentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExperimentsRequest;
  static deserializeBinaryFromReader(message: ListExperimentsRequest, reader: jspb.BinaryReader): ListExperimentsRequest;
}

export namespace ListExperimentsRequest {
  export type AsObject = {
    featureId: string,
    featureVersion?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    from: number,
    to: number,
    pageSize: number,
    cursor: string,
    environmentNamespace: string,
    status?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    maintainer: string,
    orderBy: ListExperimentsRequest.OrderByMap[keyof ListExperimentsRequest.OrderByMap],
    orderDirection: ListExperimentsRequest.OrderDirectionMap[keyof ListExperimentsRequest.OrderDirectionMap],
    searchKeyword: string,
    archived?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    statusesList: Array<proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap]>,
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

export class ListExperimentsResponse extends jspb.Message {
  clearExperimentsList(): void;
  getExperimentsList(): Array<proto_experiment_experiment_pb.Experiment>;
  setExperimentsList(value: Array<proto_experiment_experiment_pb.Experiment>): void;
  addExperiments(value?: proto_experiment_experiment_pb.Experiment, index?: number): proto_experiment_experiment_pb.Experiment;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExperimentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListExperimentsResponse): ListExperimentsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListExperimentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExperimentsResponse;
  static deserializeBinaryFromReader(message: ListExperimentsResponse, reader: jspb.BinaryReader): ListExperimentsResponse;
}

export namespace ListExperimentsResponse {
  export type AsObject = {
    experimentsList: Array<proto_experiment_experiment_pb.Experiment.AsObject>,
    cursor: string,
    totalCount: number,
  }
}

export class CreateExperimentRequest extends jspb.Message {
  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.CreateExperimentCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.CreateExperimentCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateExperimentRequest): CreateExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateExperimentRequest;
  static deserializeBinaryFromReader(message: CreateExperimentRequest, reader: jspb.BinaryReader): CreateExperimentRequest;
}

export namespace CreateExperimentRequest {
  export type AsObject = {
    command?: proto_experiment_command_pb.CreateExperimentCommand.AsObject,
    environmentNamespace: string,
  }
}

export class CreateExperimentResponse extends jspb.Message {
  hasExperiment(): boolean;
  clearExperiment(): void;
  getExperiment(): proto_experiment_experiment_pb.Experiment | undefined;
  setExperiment(value?: proto_experiment_experiment_pb.Experiment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateExperimentResponse): CreateExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateExperimentResponse;
  static deserializeBinaryFromReader(message: CreateExperimentResponse, reader: jspb.BinaryReader): CreateExperimentResponse;
}

export namespace CreateExperimentResponse {
  export type AsObject = {
    experiment?: proto_experiment_experiment_pb.Experiment.AsObject,
  }
}

export class UpdateExperimentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  hasChangeExperimentPeriodCommand(): boolean;
  clearChangeExperimentPeriodCommand(): void;
  getChangeExperimentPeriodCommand(): proto_experiment_command_pb.ChangeExperimentPeriodCommand | undefined;
  setChangeExperimentPeriodCommand(value?: proto_experiment_command_pb.ChangeExperimentPeriodCommand): void;

  hasChangeNameCommand(): boolean;
  clearChangeNameCommand(): void;
  getChangeNameCommand(): proto_experiment_command_pb.ChangeExperimentNameCommand | undefined;
  setChangeNameCommand(value?: proto_experiment_command_pb.ChangeExperimentNameCommand): void;

  hasChangeDescriptionCommand(): boolean;
  clearChangeDescriptionCommand(): void;
  getChangeDescriptionCommand(): proto_experiment_command_pb.ChangeExperimentDescriptionCommand | undefined;
  setChangeDescriptionCommand(value?: proto_experiment_command_pb.ChangeExperimentDescriptionCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateExperimentRequest): UpdateExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateExperimentRequest;
  static deserializeBinaryFromReader(message: UpdateExperimentRequest, reader: jspb.BinaryReader): UpdateExperimentRequest;
}

export namespace UpdateExperimentRequest {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    changeExperimentPeriodCommand?: proto_experiment_command_pb.ChangeExperimentPeriodCommand.AsObject,
    changeNameCommand?: proto_experiment_command_pb.ChangeExperimentNameCommand.AsObject,
    changeDescriptionCommand?: proto_experiment_command_pb.ChangeExperimentDescriptionCommand.AsObject,
  }
}

export class UpdateExperimentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateExperimentResponse): UpdateExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateExperimentResponse;
  static deserializeBinaryFromReader(message: UpdateExperimentResponse, reader: jspb.BinaryReader): UpdateExperimentResponse;
}

export namespace UpdateExperimentResponse {
  export type AsObject = {
  }
}

export class StartExperimentRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.StartExperimentCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.StartExperimentCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartExperimentRequest): StartExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartExperimentRequest;
  static deserializeBinaryFromReader(message: StartExperimentRequest, reader: jspb.BinaryReader): StartExperimentRequest;
}

export namespace StartExperimentRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_experiment_command_pb.StartExperimentCommand.AsObject,
  }
}

export class StartExperimentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartExperimentResponse): StartExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartExperimentResponse;
  static deserializeBinaryFromReader(message: StartExperimentResponse, reader: jspb.BinaryReader): StartExperimentResponse;
}

export namespace StartExperimentResponse {
  export type AsObject = {
  }
}

export class FinishExperimentRequest extends jspb.Message {
  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.FinishExperimentCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.FinishExperimentCommand): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FinishExperimentRequest): FinishExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FinishExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishExperimentRequest;
  static deserializeBinaryFromReader(message: FinishExperimentRequest, reader: jspb.BinaryReader): FinishExperimentRequest;
}

export namespace FinishExperimentRequest {
  export type AsObject = {
    environmentNamespace: string,
    id: string,
    command?: proto_experiment_command_pb.FinishExperimentCommand.AsObject,
  }
}

export class FinishExperimentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FinishExperimentResponse): FinishExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FinishExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishExperimentResponse;
  static deserializeBinaryFromReader(message: FinishExperimentResponse, reader: jspb.BinaryReader): FinishExperimentResponse;
}

export namespace FinishExperimentResponse {
  export type AsObject = {
  }
}

export class StopExperimentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.StopExperimentCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.StopExperimentCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopExperimentRequest): StopExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopExperimentRequest;
  static deserializeBinaryFromReader(message: StopExperimentRequest, reader: jspb.BinaryReader): StopExperimentRequest;
}

export namespace StopExperimentRequest {
  export type AsObject = {
    id: string,
    command?: proto_experiment_command_pb.StopExperimentCommand.AsObject,
    environmentNamespace: string,
  }
}

export class StopExperimentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StopExperimentResponse): StopExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopExperimentResponse;
  static deserializeBinaryFromReader(message: StopExperimentResponse, reader: jspb.BinaryReader): StopExperimentResponse;
}

export namespace StopExperimentResponse {
  export type AsObject = {
  }
}

export class ArchiveExperimentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.ArchiveExperimentCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.ArchiveExperimentCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveExperimentRequest): ArchiveExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveExperimentRequest;
  static deserializeBinaryFromReader(message: ArchiveExperimentRequest, reader: jspb.BinaryReader): ArchiveExperimentRequest;
}

export namespace ArchiveExperimentRequest {
  export type AsObject = {
    id: string,
    command?: proto_experiment_command_pb.ArchiveExperimentCommand.AsObject,
    environmentNamespace: string,
  }
}

export class ArchiveExperimentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveExperimentResponse): ArchiveExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveExperimentResponse;
  static deserializeBinaryFromReader(message: ArchiveExperimentResponse, reader: jspb.BinaryReader): ArchiveExperimentResponse;
}

export namespace ArchiveExperimentResponse {
  export type AsObject = {
  }
}

export class DeleteExperimentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): proto_experiment_command_pb.DeleteExperimentCommand | undefined;
  setCommand(value?: proto_experiment_command_pb.DeleteExperimentCommand): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteExperimentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteExperimentRequest): DeleteExperimentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteExperimentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteExperimentRequest;
  static deserializeBinaryFromReader(message: DeleteExperimentRequest, reader: jspb.BinaryReader): DeleteExperimentRequest;
}

export namespace DeleteExperimentRequest {
  export type AsObject = {
    id: string,
    command?: proto_experiment_command_pb.DeleteExperimentCommand.AsObject,
    environmentNamespace: string,
  }
}

export class DeleteExperimentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteExperimentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteExperimentResponse): DeleteExperimentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteExperimentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteExperimentResponse;
  static deserializeBinaryFromReader(message: DeleteExperimentResponse, reader: jspb.BinaryReader): DeleteExperimentResponse;
}

export namespace DeleteExperimentResponse {
  export type AsObject = {
  }
}

