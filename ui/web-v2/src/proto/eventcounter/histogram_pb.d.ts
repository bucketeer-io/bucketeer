// package: bucketeer.eventcounter
// file: proto/eventcounter/histogram.proto

import * as jspb from "google-protobuf";

export class Histogram extends jspb.Message {
  clearHistList(): void;
  getHistList(): Array<number>;
  setHistList(value: Array<number>): void;
  addHist(value: number, index?: number): number;

  clearBinsList(): void;
  getBinsList(): Array<number>;
  setBinsList(value: Array<number>): void;
  addBins(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Histogram.AsObject;
  static toObject(includeInstance: boolean, msg: Histogram): Histogram.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Histogram, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Histogram;
  static deserializeBinaryFromReader(message: Histogram, reader: jspb.BinaryReader): Histogram;
}

export namespace Histogram {
  export type AsObject = {
    histList: Array<number>,
    binsList: Array<number>,
  }
}

