// package: bucketeer.experiment
// file: proto/experiment/service.proto

import * as proto_experiment_service_pb from "../../proto/experiment/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ExperimentServiceGetGoal = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.GetGoalRequest;
  readonly responseType: typeof proto_experiment_service_pb.GetGoalResponse;
};

type ExperimentServiceListGoals = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.ListGoalsRequest;
  readonly responseType: typeof proto_experiment_service_pb.ListGoalsResponse;
};

type ExperimentServiceCreateGoal = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.CreateGoalRequest;
  readonly responseType: typeof proto_experiment_service_pb.CreateGoalResponse;
};

type ExperimentServiceUpdateGoal = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.UpdateGoalRequest;
  readonly responseType: typeof proto_experiment_service_pb.UpdateGoalResponse;
};

type ExperimentServiceArchiveGoal = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.ArchiveGoalRequest;
  readonly responseType: typeof proto_experiment_service_pb.ArchiveGoalResponse;
};

type ExperimentServiceDeleteGoal = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.DeleteGoalRequest;
  readonly responseType: typeof proto_experiment_service_pb.DeleteGoalResponse;
};

type ExperimentServiceGetExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.GetExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.GetExperimentResponse;
};

type ExperimentServiceListExperiments = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.ListExperimentsRequest;
  readonly responseType: typeof proto_experiment_service_pb.ListExperimentsResponse;
};

type ExperimentServiceCreateExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.CreateExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.CreateExperimentResponse;
};

type ExperimentServiceUpdateExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.UpdateExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.UpdateExperimentResponse;
};

type ExperimentServiceStartExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.StartExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.StartExperimentResponse;
};

type ExperimentServiceFinishExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.FinishExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.FinishExperimentResponse;
};

type ExperimentServiceStopExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.StopExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.StopExperimentResponse;
};

type ExperimentServiceArchiveExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.ArchiveExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.ArchiveExperimentResponse;
};

type ExperimentServiceDeleteExperiment = {
  readonly methodName: string;
  readonly service: typeof ExperimentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_experiment_service_pb.DeleteExperimentRequest;
  readonly responseType: typeof proto_experiment_service_pb.DeleteExperimentResponse;
};

export class ExperimentService {
  static readonly serviceName: string;
  static readonly GetGoal: ExperimentServiceGetGoal;
  static readonly ListGoals: ExperimentServiceListGoals;
  static readonly CreateGoal: ExperimentServiceCreateGoal;
  static readonly UpdateGoal: ExperimentServiceUpdateGoal;
  static readonly ArchiveGoal: ExperimentServiceArchiveGoal;
  static readonly DeleteGoal: ExperimentServiceDeleteGoal;
  static readonly GetExperiment: ExperimentServiceGetExperiment;
  static readonly ListExperiments: ExperimentServiceListExperiments;
  static readonly CreateExperiment: ExperimentServiceCreateExperiment;
  static readonly UpdateExperiment: ExperimentServiceUpdateExperiment;
  static readonly StartExperiment: ExperimentServiceStartExperiment;
  static readonly FinishExperiment: ExperimentServiceFinishExperiment;
  static readonly StopExperiment: ExperimentServiceStopExperiment;
  static readonly ArchiveExperiment: ExperimentServiceArchiveExperiment;
  static readonly DeleteExperiment: ExperimentServiceDeleteExperiment;
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

export class ExperimentServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getGoal(
    requestMessage: proto_experiment_service_pb.GetGoalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.GetGoalResponse|null) => void
  ): UnaryResponse;
  getGoal(
    requestMessage: proto_experiment_service_pb.GetGoalRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.GetGoalResponse|null) => void
  ): UnaryResponse;
  listGoals(
    requestMessage: proto_experiment_service_pb.ListGoalsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ListGoalsResponse|null) => void
  ): UnaryResponse;
  listGoals(
    requestMessage: proto_experiment_service_pb.ListGoalsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ListGoalsResponse|null) => void
  ): UnaryResponse;
  createGoal(
    requestMessage: proto_experiment_service_pb.CreateGoalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.CreateGoalResponse|null) => void
  ): UnaryResponse;
  createGoal(
    requestMessage: proto_experiment_service_pb.CreateGoalRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.CreateGoalResponse|null) => void
  ): UnaryResponse;
  updateGoal(
    requestMessage: proto_experiment_service_pb.UpdateGoalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.UpdateGoalResponse|null) => void
  ): UnaryResponse;
  updateGoal(
    requestMessage: proto_experiment_service_pb.UpdateGoalRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.UpdateGoalResponse|null) => void
  ): UnaryResponse;
  archiveGoal(
    requestMessage: proto_experiment_service_pb.ArchiveGoalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ArchiveGoalResponse|null) => void
  ): UnaryResponse;
  archiveGoal(
    requestMessage: proto_experiment_service_pb.ArchiveGoalRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ArchiveGoalResponse|null) => void
  ): UnaryResponse;
  deleteGoal(
    requestMessage: proto_experiment_service_pb.DeleteGoalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.DeleteGoalResponse|null) => void
  ): UnaryResponse;
  deleteGoal(
    requestMessage: proto_experiment_service_pb.DeleteGoalRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.DeleteGoalResponse|null) => void
  ): UnaryResponse;
  getExperiment(
    requestMessage: proto_experiment_service_pb.GetExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.GetExperimentResponse|null) => void
  ): UnaryResponse;
  getExperiment(
    requestMessage: proto_experiment_service_pb.GetExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.GetExperimentResponse|null) => void
  ): UnaryResponse;
  listExperiments(
    requestMessage: proto_experiment_service_pb.ListExperimentsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ListExperimentsResponse|null) => void
  ): UnaryResponse;
  listExperiments(
    requestMessage: proto_experiment_service_pb.ListExperimentsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ListExperimentsResponse|null) => void
  ): UnaryResponse;
  createExperiment(
    requestMessage: proto_experiment_service_pb.CreateExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.CreateExperimentResponse|null) => void
  ): UnaryResponse;
  createExperiment(
    requestMessage: proto_experiment_service_pb.CreateExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.CreateExperimentResponse|null) => void
  ): UnaryResponse;
  updateExperiment(
    requestMessage: proto_experiment_service_pb.UpdateExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.UpdateExperimentResponse|null) => void
  ): UnaryResponse;
  updateExperiment(
    requestMessage: proto_experiment_service_pb.UpdateExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.UpdateExperimentResponse|null) => void
  ): UnaryResponse;
  startExperiment(
    requestMessage: proto_experiment_service_pb.StartExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.StartExperimentResponse|null) => void
  ): UnaryResponse;
  startExperiment(
    requestMessage: proto_experiment_service_pb.StartExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.StartExperimentResponse|null) => void
  ): UnaryResponse;
  finishExperiment(
    requestMessage: proto_experiment_service_pb.FinishExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.FinishExperimentResponse|null) => void
  ): UnaryResponse;
  finishExperiment(
    requestMessage: proto_experiment_service_pb.FinishExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.FinishExperimentResponse|null) => void
  ): UnaryResponse;
  stopExperiment(
    requestMessage: proto_experiment_service_pb.StopExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.StopExperimentResponse|null) => void
  ): UnaryResponse;
  stopExperiment(
    requestMessage: proto_experiment_service_pb.StopExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.StopExperimentResponse|null) => void
  ): UnaryResponse;
  archiveExperiment(
    requestMessage: proto_experiment_service_pb.ArchiveExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ArchiveExperimentResponse|null) => void
  ): UnaryResponse;
  archiveExperiment(
    requestMessage: proto_experiment_service_pb.ArchiveExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.ArchiveExperimentResponse|null) => void
  ): UnaryResponse;
  deleteExperiment(
    requestMessage: proto_experiment_service_pb.DeleteExperimentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.DeleteExperimentResponse|null) => void
  ): UnaryResponse;
  deleteExperiment(
    requestMessage: proto_experiment_service_pb.DeleteExperimentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_experiment_service_pb.DeleteExperimentResponse|null) => void
  ): UnaryResponse;
}

