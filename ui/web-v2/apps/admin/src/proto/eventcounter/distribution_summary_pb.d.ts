// package: bucketeer.eventcounter
// file: proto/eventcounter/distribution_summary.proto

import * as jspb from "google-protobuf";
import * as proto_eventcounter_histogram_pb from "../../proto/eventcounter/histogram_pb";

export class DistributionSummary extends jspb.Message {
  getMean(): number;
  setMean(value: number): void;

  getSd(): number;
  setSd(value: number): void;

  getRhat(): number;
  setRhat(value: number): void;

  hasHistogram(): boolean;
  clearHistogram(): void;
  getHistogram(): proto_eventcounter_histogram_pb.Histogram | undefined;
  setHistogram(value?: proto_eventcounter_histogram_pb.Histogram): void;

  getMedian(): number;
  setMedian(value: number): void;

  getPercentile025(): number;
  setPercentile025(value: number): void;

  getPercentile975(): number;
  setPercentile975(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DistributionSummary.AsObject;
  static toObject(includeInstance: boolean, msg: DistributionSummary): DistributionSummary.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DistributionSummary, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DistributionSummary;
  static deserializeBinaryFromReader(message: DistributionSummary, reader: jspb.BinaryReader): DistributionSummary;
}

export namespace DistributionSummary {
  export type AsObject = {
    mean: number,
    sd: number,
    rhat: number,
    histogram?: proto_eventcounter_histogram_pb.Histogram.AsObject,
    median: number,
    percentile025: number,
    percentile975: number,
  }
}

