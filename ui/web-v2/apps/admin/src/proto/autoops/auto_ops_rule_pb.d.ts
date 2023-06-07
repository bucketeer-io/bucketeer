// package: bucketeer.autoops
// file: proto/autoops/auto_ops_rule.proto

import * as jspb from "google-protobuf";
import * as proto_autoops_clause_pb from "../../proto/autoops/clause_pb";

export class AutoOpsRule extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getOpsType(): OpsTypeMap[keyof OpsTypeMap];
  setOpsType(value: OpsTypeMap[keyof OpsTypeMap]): void;

  clearClausesList(): void;
  getClausesList(): Array<proto_autoops_clause_pb.Clause>;
  setClausesList(value: Array<proto_autoops_clause_pb.Clause>): void;
  addClauses(value?: proto_autoops_clause_pb.Clause, index?: number): proto_autoops_clause_pb.Clause;

  getTriggeredAt(): number;
  setTriggeredAt(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoOpsRule.AsObject;
  static toObject(includeInstance: boolean, msg: AutoOpsRule): AutoOpsRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AutoOpsRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRule;
  static deserializeBinaryFromReader(message: AutoOpsRule, reader: jspb.BinaryReader): AutoOpsRule;
}

export namespace AutoOpsRule {
  export type AsObject = {
    id: string,
    featureId: string,
    opsType: OpsTypeMap[keyof OpsTypeMap],
    clausesList: Array<proto_autoops_clause_pb.Clause.AsObject>,
    triggeredAt: number,
    createdAt: number,
    updatedAt: number,
    deleted: boolean,
  }
}

export interface OpsTypeMap {
  ENABLE_FEATURE: 0;
  DISABLE_FEATURE: 1;
}

export const OpsType: OpsTypeMap;

