// package: bucketeer.experimentcalculator
// file: proto/experimentcalculator/service.proto

import * as proto_experimentcalculator_service_pb from "../../proto/experimentcalculator/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ExperimentCalculatorServiceCalcExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentCalculatorService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experimentcalculator_service_pb.BatchCalcRequest;
  readonly responseType: typeof proto_experimentcalculator_service_pb.BatchCalcResponse;
};

export class ExperimentCalculatorService {
  static readonly serviceName: string;
  static readonly CalcExperiment: ExperimentCalculatorServiceCalcExperiment;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class ExperimentCalculatorServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  calcExperiment(
    requestMessage: proto_experimentcalculator_service_pb.BatchCalcRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experimentcalculator_service_pb.BatchCalcResponse|null) => void
  ): UnaryResponse;
  calcExperiment(
    requestMessage: proto_experimentcalculator_service_pb.BatchCalcRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experimentcalculator_service_pb.BatchCalcResponse|null) => void
  ): UnaryResponse;
}

