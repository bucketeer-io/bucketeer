// package: bucketeer.feature
// file: proto/feature/reason.proto

import * as jspb from "google-protobuf";

export class Reason extends jspb.Message {
  getType(): Reason.TypeMap[keyof Reason.TypeMap];
  setType(value: Reason.TypeMap[keyof Reason.TypeMap]): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Reason.AsObject;
  static toObject(includeInstance: boolean, msg: Reason): Reason.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Reason, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Reason;
  static deserializeBinaryFromReader(message: Reason, reader: jspb.BinaryReader): Reason;
}

export namespace Reason {
  export type AsObject = {
    type: Reason.TypeMap[keyof Reason.TypeMap],
    ruleId: string,
  }

  export interface TypeMap {
    TARGET: 0;
    RULE: 1;
    DEFAULT: 3;
    CLIENT: 4;
    OFF_VARIATION: 5;
    PREREQUISITE: 6;
  }

  export const Type: TypeMap;
}

