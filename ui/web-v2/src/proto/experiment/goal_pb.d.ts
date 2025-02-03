// package: bucketeer.experiment
// file: proto/experiment/goal.proto

import * as jspb from 'google-protobuf';
import * as proto_experiment_experiment_pb from '../../proto/experiment/experiment_pb';
import * as proto_autoops_auto_ops_rule_pb from '../../proto/autoops/auto_ops_rule_pb';

export class Goal extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getIsInUseStatus(): boolean;
  setIsInUseStatus(value: boolean): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  getConnectionType(): Goal.ConnectionTypeMap[keyof Goal.ConnectionTypeMap];
  setConnectionType(
    value: Goal.ConnectionTypeMap[keyof Goal.ConnectionTypeMap]
  ): void;

  clearExperimentsList(): void;
  getExperimentsList(): Array<Goal.ExperimentReference>;
  setExperimentsList(value: Array<Goal.ExperimentReference>): void;
  addExperiments(
    value?: Goal.ExperimentReference,
    index?: number
  ): Goal.ExperimentReference;

  clearAutoOpsRulesList(): void;
  getAutoOpsRulesList(): Array<Goal.AutoOpsRuleReference>;
  setAutoOpsRulesList(value: Array<Goal.AutoOpsRuleReference>): void;
  addAutoOpsRules(
    value?: Goal.AutoOpsRuleReference,
    index?: number
  ): Goal.AutoOpsRuleReference;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Goal.AsObject;
  static toObject(includeInstance: boolean, msg: Goal): Goal.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Goal,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Goal;
  static deserializeBinaryFromReader(
    message: Goal,
    reader: jspb.BinaryReader
  ): Goal;
}

export namespace Goal {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    deleted: boolean;
    createdAt: number;
    updatedAt: number;
    isInUseStatus: boolean;
    archived: boolean;
    connectionType: Goal.ConnectionTypeMap[keyof Goal.ConnectionTypeMap];
    experimentsList: Array<Goal.ExperimentReference.AsObject>;
    autoOpsRulesList: Array<Goal.AutoOpsRuleReference.AsObject>;
  };

  export class ExperimentReference extends jspb.Message {
    getId(): string;
    setId(value: string): void;

    getName(): string;
    setName(value: string): void;

    getFeatureId(): string;
    setFeatureId(value: string): void;

    getStatus(): proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap];
    setStatus(
      value: proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap]
    ): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExperimentReference.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: ExperimentReference
    ): ExperimentReference.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: ExperimentReference,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): ExperimentReference;
    static deserializeBinaryFromReader(
      message: ExperimentReference,
      reader: jspb.BinaryReader
    ): ExperimentReference;
  }

  export namespace ExperimentReference {
    export type AsObject = {
      id: string;
      name: string;
      featureId: string;
      status: proto_experiment_experiment_pb.Experiment.StatusMap[keyof proto_experiment_experiment_pb.Experiment.StatusMap];
    };
  }

  export class AutoOpsRuleReference extends jspb.Message {
    getId(): string;
    setId(value: string): void;

    getFeatureId(): string;
    setFeatureId(value: string): void;

    getAutoOpsStatus(): proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
    setAutoOpsStatus(
      value: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap]
    ): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AutoOpsRuleReference.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: AutoOpsRuleReference
    ): AutoOpsRuleReference.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: AutoOpsRuleReference,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): AutoOpsRuleReference;
    static deserializeBinaryFromReader(
      message: AutoOpsRuleReference,
      reader: jspb.BinaryReader
    ): AutoOpsRuleReference;
  }

  export namespace AutoOpsRuleReference {
    export type AsObject = {
      id: string;
      featureId: string;
      autoOpsStatus: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
    };
  }

  export interface ConnectionTypeMap {
    UNKNOWN: 0;
    EXPERIMENT: 1;
    OPERATION: 2;
  }

  export const ConnectionType: ConnectionTypeMap;
}
