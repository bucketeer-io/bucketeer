// package: bucketeer.feature
// file: proto/feature/clause.proto

import * as jspb from "google-protobuf";

export class Clause extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAttribute(): string;
  setAttribute(value: string): void;

  getOperator(): Clause.OperatorMap[keyof Clause.OperatorMap];
  setOperator(value: Clause.OperatorMap[keyof Clause.OperatorMap]): void;

  clearValuesList(): void;
  getValuesList(): Array<string>;
  setValuesList(value: Array<string>): void;
  addValues(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Clause.AsObject;
  static toObject(includeInstance: boolean, msg: Clause): Clause.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Clause, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Clause;
  static deserializeBinaryFromReader(message: Clause, reader: jspb.BinaryReader): Clause;
}

export namespace Clause {
  export type AsObject = {
    id: string,
    attribute: string,
    operator: Clause.OperatorMap[keyof Clause.OperatorMap],
    valuesList: Array<string>,
  }

  export interface OperatorMap {
    EQUALS: 0;
    IN: 1;
    ENDS_WITH: 2;
    STARTS_WITH: 3;
    SEGMENT: 4;
    GREATER: 5;
    GREATER_OR_EQUAL: 6;
    LESS: 7;
    LESS_OR_EQUAL: 8;
    BEFORE: 9;
    AFTER: 10;
  }

  export const Operator: OperatorMap;
}

