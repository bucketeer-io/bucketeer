// package: bucketeer.experimentcalculator
// file: proto/experimentcalculator/service.proto

import * as jspb from "google-protobuf";
import * as proto_experiment_experiment_pb from "../../proto/experiment/experiment_pb";

export class BatchCalcRequest extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  hasExperiment(): boolean;
  clearExperiment(): void;
  getExperiment(): proto_experiment_experiment_pb.Experiment | undefined;
  setExperiment(value?: proto_experiment_experiment_pb.Experiment): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchCalcRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BatchCalcRequest): BatchCalcRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BatchCalcRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchCalcRequest;
  static deserializeBinaryFromReader(message: BatchCalcRequest, reader: jspb.BinaryReader): BatchCalcRequest;
}

export namespace BatchCalcRequest {
  export type AsObject = {
    environmentId: string,
    experiment?: proto_experiment_experiment_pb.Experiment.AsObject,
  }
}

export class BatchCalcResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchCalcResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BatchCalcResponse): BatchCalcResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BatchCalcResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchCalcResponse;
  static deserializeBinaryFromReader(message: BatchCalcResponse, reader: jspb.BinaryReader): BatchCalcResponse;
}

export namespace BatchCalcResponse {
  export type AsObject = {
  }
}

