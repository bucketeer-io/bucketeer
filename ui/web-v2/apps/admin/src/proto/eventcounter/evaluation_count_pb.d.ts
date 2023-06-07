// package: bucketeer.eventcounter
// file: proto/eventcounter/evaluation_count.proto

import * as jspb from "google-protobuf";
import * as proto_eventcounter_variation_count_pb from "../../proto/eventcounter/variation_count_pb";

export class EvaluationCount extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  clearRealtimeCountsList(): void;
  getRealtimeCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setRealtimeCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addRealtimeCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  clearBatchCountsList(): void;
  getBatchCountsList(): Array<proto_eventcounter_variation_count_pb.VariationCount>;
  setBatchCountsList(value: Array<proto_eventcounter_variation_count_pb.VariationCount>): void;
  addBatchCounts(value?: proto_eventcounter_variation_count_pb.VariationCount, index?: number): proto_eventcounter_variation_count_pb.VariationCount;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluationCount.AsObject;
  static toObject(includeInstance: boolean, msg: EvaluationCount): EvaluationCount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EvaluationCount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EvaluationCount;
  static deserializeBinaryFromReader(message: EvaluationCount, reader: jspb.BinaryReader): EvaluationCount;
}

export namespace EvaluationCount {
  export type AsObject = {
    id: string,
    featureId: string,
    featureVersion: number,
    realtimeCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
    batchCountsList: Array<proto_eventcounter_variation_count_pb.VariationCount.AsObject>,
    updatedAt: number,
  }
}

