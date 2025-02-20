// package: bucketeer.common
// file: proto/common/string.proto

import * as jspb from 'google-protobuf';

export class StringListValue extends jspb.Message {
  clearValuesList(): void;
  getValuesList(): Array<string>;
  setValuesList(value: Array<string>): void;
  addValues(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StringListValue.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StringListValue
  ): StringListValue.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StringListValue,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StringListValue;
  static deserializeBinaryFromReader(
    message: StringListValue,
    reader: jspb.BinaryReader
  ): StringListValue;
}

export namespace StringListValue {
  export type AsObject = {
    valuesList: Array<string>;
  };
}
