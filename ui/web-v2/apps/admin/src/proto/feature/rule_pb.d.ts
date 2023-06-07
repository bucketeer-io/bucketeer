// package: bucketeer.feature
// file: proto/feature/rule.proto

import * as jspb from "google-protobuf";
import * as proto_feature_clause_pb from "../../proto/feature/clause_pb";
import * as proto_feature_strategy_pb from "../../proto/feature/strategy_pb";

export class Rule extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasStrategy(): boolean;
  clearStrategy(): void;
  getStrategy(): proto_feature_strategy_pb.Strategy | undefined;
  setStrategy(value?: proto_feature_strategy_pb.Strategy): void;

  clearClausesList(): void;
  getClausesList(): Array<proto_feature_clause_pb.Clause>;
  setClausesList(value: Array<proto_feature_clause_pb.Clause>): void;
  addClauses(value?: proto_feature_clause_pb.Clause, index?: number): proto_feature_clause_pb.Clause;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Rule.AsObject;
  static toObject(includeInstance: boolean, msg: Rule): Rule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Rule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Rule;
  static deserializeBinaryFromReader(message: Rule, reader: jspb.BinaryReader): Rule;
}

export namespace Rule {
  export type AsObject = {
    id: string,
    strategy?: proto_feature_strategy_pb.Strategy.AsObject,
    clausesList: Array<proto_feature_clause_pb.Clause.AsObject>,
  }
}

