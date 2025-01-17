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
  getExperimentsList(): Array<proto_experiment_experiment_pb.Experiment>;
  setExperimentsList(
    value: Array<proto_experiment_experiment_pb.Experiment>
  ): void;
  addExperiments(
    value?: proto_experiment_experiment_pb.Experiment,
    index?: number
  ): proto_experiment_experiment_pb.Experiment;

  clearAutoOpsRulesList(): void;
  getAutoOpsRulesList(): Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule>;
  setAutoOpsRulesList(
    value: Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule>
  ): void;
  addAutoOpsRules(
    value?: proto_autoops_auto_ops_rule_pb.AutoOpsRule,
    index?: number
  ): proto_autoops_auto_ops_rule_pb.AutoOpsRule;

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
    experimentsList: Array<proto_experiment_experiment_pb.Experiment.AsObject>;
    autoOpsRulesList: Array<proto_autoops_auto_ops_rule_pb.AutoOpsRule.AsObject>;
  };

  export interface ConnectionTypeMap {
    UNKNOWN: 0;
    EXPERIMENT: 1;
    OPERATION: 2;
  }

  export const ConnectionType: ConnectionTypeMap;
}
