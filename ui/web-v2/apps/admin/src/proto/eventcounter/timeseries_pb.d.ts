// package: bucketeer.eventcounter
// file: proto/eventcounter/timeseries.proto

import * as jspb from "google-protobuf";

export class VariationTimeseries extends jspb.Message {
  getVariationId(): string;
  setVariationId(value: string): void;

  hasTimeseries(): boolean;
  clearTimeseries(): void;
  getTimeseries(): Timeseries | undefined;
  setTimeseries(value?: Timeseries): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationTimeseries.AsObject;
  static toObject(includeInstance: boolean, msg: VariationTimeseries): VariationTimeseries.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VariationTimeseries, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VariationTimeseries;
  static deserializeBinaryFromReader(message: VariationTimeseries, reader: jspb.BinaryReader): VariationTimeseries;
}

export namespace VariationTimeseries {
  export type AsObject = {
    variationId: string,
    timeseries?: Timeseries.AsObject,
  }
}

export class Timeseries extends jspb.Message {
  clearTimestampsList(): void;
  getTimestampsList(): Array<number>;
  setTimestampsList(value: Array<number>): void;
  addTimestamps(value: number, index?: number): number;

  clearValuesList(): void;
  getValuesList(): Array<number>;
  setValuesList(value: Array<number>): void;
  addValues(value: number, index?: number): number;

  getUnit(): Timeseries.UnitMap[keyof Timeseries.UnitMap];
  setUnit(value: Timeseries.UnitMap[keyof Timeseries.UnitMap]): void;

  getTotalCounts(): number;
  setTotalCounts(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Timeseries.AsObject;
  static toObject(includeInstance: boolean, msg: Timeseries): Timeseries.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Timeseries, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Timeseries;
  static deserializeBinaryFromReader(message: Timeseries, reader: jspb.BinaryReader): Timeseries;
}

export namespace Timeseries {
  export type AsObject = {
    timestampsList: Array<number>,
    valuesList: Array<number>,
    unit: Timeseries.UnitMap[keyof Timeseries.UnitMap],
    totalCounts: number,
  }

  export interface UnitMap {
    HOUR: 0;
    DAY: 1;
  }

  export const Unit: UnitMap;
}

