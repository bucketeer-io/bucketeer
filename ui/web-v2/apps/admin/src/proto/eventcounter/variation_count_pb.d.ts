// package: bucketeer.eventcounter
// file: proto/eventcounter/variation_count.proto

import * as jspb from "google-protobuf";

export class VariationCount extends jspb.Message {
  getVariationId(): string;
  setVariationId(value: string): void;

  getUserCount(): number;
  setUserCount(value: number): void;

  getEventCount(): number;
  setEventCount(value: number): void;

  getValueSum(): number;
  setValueSum(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getVariationValue(): string;
  setVariationValue(value: string): void;

  getValueSumPerUserMean(): number;
  setValueSumPerUserMean(value: number): void;

  getValueSumPerUserVariance(): number;
  setValueSumPerUserVariance(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationCount.AsObject;
  static toObject(includeInstance: boolean, msg: VariationCount): VariationCount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VariationCount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VariationCount;
  static deserializeBinaryFromReader(message: VariationCount, reader: jspb.BinaryReader): VariationCount;
}

export namespace VariationCount {
  export type AsObject = {
    variationId: string,
    userCount: number,
    eventCount: number,
    valueSum: number,
    createdAt: number,
    variationValue: string,
    valueSumPerUserMean: number,
    valueSumPerUserVariance: number,
  }
}

