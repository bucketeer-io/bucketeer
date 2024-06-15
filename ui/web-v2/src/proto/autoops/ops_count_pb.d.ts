// package: bucketeer.autoops
// file: proto/autoops/ops_count.proto

import * as jspb from "google-protobuf";

export class OpsCount extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAutoOpsRuleId(): string;
  setAutoOpsRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getOpsEventCount(): number;
  setOpsEventCount(value: number): void;

  getEvaluationCount(): number;
  setEvaluationCount(value: number): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpsCount.AsObject;
  static toObject(includeInstance: boolean, msg: OpsCount): OpsCount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpsCount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpsCount;
  static deserializeBinaryFromReader(message: OpsCount, reader: jspb.BinaryReader): OpsCount;
}

export namespace OpsCount {
  export type AsObject = {
    id: string,
    autoOpsRuleId: string,
    clauseId: string,
    updatedAt: number,
    opsEventCount: number,
    evaluationCount: number,
    featureId: string,
  }
}

