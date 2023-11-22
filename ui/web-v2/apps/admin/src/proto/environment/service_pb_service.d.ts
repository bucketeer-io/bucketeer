// package: bucketeer.environment
// file: proto/environment/service.proto

import * as proto_environment_service_pb from "../../proto/environment/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type EnvironmentServiceGetEnvironmentV2 = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.GetEnvironmentV2Request;
  readonly responseType: typeof proto_environment_service_pb.GetEnvironmentV2Response;
};

type EnvironmentServiceListEnvironmentsV2 = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ListEnvironmentsV2Request;
  readonly responseType: typeof proto_environment_service_pb.ListEnvironmentsV2Response;
};

type EnvironmentServiceCreateEnvironmentV2 = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.CreateEnvironmentV2Request;
  readonly responseType: typeof proto_environment_service_pb.CreateEnvironmentV2Response;
};

type EnvironmentServiceUpdateEnvironmentV2 = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.UpdateEnvironmentV2Request;
  readonly responseType: typeof proto_environment_service_pb.UpdateEnvironmentV2Response;
};

type EnvironmentServiceArchiveEnvironmentV2 = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ArchiveEnvironmentV2Request;
  readonly responseType: typeof proto_environment_service_pb.ArchiveEnvironmentV2Response;
};

type EnvironmentServiceUnarchiveEnvironmentV2 = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.UnarchiveEnvironmentV2Request;
  readonly responseType: typeof proto_environment_service_pb.UnarchiveEnvironmentV2Response;
};

type EnvironmentServiceGetProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.GetProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.GetProjectResponse;
};

type EnvironmentServiceListProjects = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ListProjectsRequest;
  readonly responseType: typeof proto_environment_service_pb.ListProjectsResponse;
};

type EnvironmentServiceCreateProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.CreateProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.CreateProjectResponse;
};

type EnvironmentServiceCreateTrialProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.CreateTrialProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.CreateTrialProjectResponse;
};

type EnvironmentServiceUpdateProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.UpdateProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.UpdateProjectResponse;
};

type EnvironmentServiceEnableProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.EnableProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.EnableProjectResponse;
};

type EnvironmentServiceDisableProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.DisableProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.DisableProjectResponse;
};

type EnvironmentServiceConvertTrialProject = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ConvertTrialProjectRequest;
  readonly responseType: typeof proto_environment_service_pb.ConvertTrialProjectResponse;
};

type EnvironmentServiceGetOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.GetOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.GetOrganizationResponse;
};

type EnvironmentServiceListOrganizations = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ListOrganizationsRequest;
  readonly responseType: typeof proto_environment_service_pb.ListOrganizationsResponse;
};

type EnvironmentServiceCreateOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.CreateOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.CreateOrganizationResponse;
};

type EnvironmentServiceUpdateOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.UpdateOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.UpdateOrganizationResponse;
};

type EnvironmentServiceEnableOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.EnableOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.EnableOrganizationResponse;
};

type EnvironmentServiceDisableOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.DisableOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.DisableOrganizationResponse;
};

type EnvironmentServiceArchiveOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ArchiveOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.ArchiveOrganizationResponse;
};

type EnvironmentServiceUnarchiveOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.UnarchiveOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.UnarchiveOrganizationResponse;
};

type EnvironmentServiceConvertTrialOrganization = {
  readonly methodName: string;
  readonly service: typeof EnvironmentService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_environment_service_pb.ConvertTrialOrganizationRequest;
  readonly responseType: typeof proto_environment_service_pb.ConvertTrialOrganizationResponse;
};

export class EnvironmentService {
  static readonly serviceName: string;
  static readonly GetEnvironmentV2: EnvironmentServiceGetEnvironmentV2;
  static readonly ListEnvironmentsV2: EnvironmentServiceListEnvironmentsV2;
  static readonly CreateEnvironmentV2: EnvironmentServiceCreateEnvironmentV2;
  static readonly UpdateEnvironmentV2: EnvironmentServiceUpdateEnvironmentV2;
  static readonly ArchiveEnvironmentV2: EnvironmentServiceArchiveEnvironmentV2;
  static readonly UnarchiveEnvironmentV2: EnvironmentServiceUnarchiveEnvironmentV2;
  static readonly GetProject: EnvironmentServiceGetProject;
  static readonly ListProjects: EnvironmentServiceListProjects;
  static readonly CreateProject: EnvironmentServiceCreateProject;
  static readonly CreateTrialProject: EnvironmentServiceCreateTrialProject;
  static readonly UpdateProject: EnvironmentServiceUpdateProject;
  static readonly EnableProject: EnvironmentServiceEnableProject;
  static readonly DisableProject: EnvironmentServiceDisableProject;
  static readonly ConvertTrialProject: EnvironmentServiceConvertTrialProject;
  static readonly GetOrganization: EnvironmentServiceGetOrganization;
  static readonly ListOrganizations: EnvironmentServiceListOrganizations;
  static readonly CreateOrganization: EnvironmentServiceCreateOrganization;
  static readonly UpdateOrganization: EnvironmentServiceUpdateOrganization;
  static readonly EnableOrganization: EnvironmentServiceEnableOrganization;
  static readonly DisableOrganization: EnvironmentServiceDisableOrganization;
  static readonly ArchiveOrganization: EnvironmentServiceArchiveOrganization;
  static readonly UnarchiveOrganization: EnvironmentServiceUnarchiveOrganization;
  static readonly ConvertTrialOrganization: EnvironmentServiceConvertTrialOrganization;
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

export class EnvironmentServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getEnvironmentV2(
    requestMessage: proto_environment_service_pb.GetEnvironmentV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.GetEnvironmentV2Response|null) => void
  ): UnaryResponse;
  getEnvironmentV2(
    requestMessage: proto_environment_service_pb.GetEnvironmentV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.GetEnvironmentV2Response|null) => void
  ): UnaryResponse;
  listEnvironmentsV2(
    requestMessage: proto_environment_service_pb.ListEnvironmentsV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ListEnvironmentsV2Response|null) => void
  ): UnaryResponse;
  listEnvironmentsV2(
    requestMessage: proto_environment_service_pb.ListEnvironmentsV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ListEnvironmentsV2Response|null) => void
  ): UnaryResponse;
  createEnvironmentV2(
    requestMessage: proto_environment_service_pb.CreateEnvironmentV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateEnvironmentV2Response|null) => void
  ): UnaryResponse;
  createEnvironmentV2(
    requestMessage: proto_environment_service_pb.CreateEnvironmentV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateEnvironmentV2Response|null) => void
  ): UnaryResponse;
  updateEnvironmentV2(
    requestMessage: proto_environment_service_pb.UpdateEnvironmentV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UpdateEnvironmentV2Response|null) => void
  ): UnaryResponse;
  updateEnvironmentV2(
    requestMessage: proto_environment_service_pb.UpdateEnvironmentV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UpdateEnvironmentV2Response|null) => void
  ): UnaryResponse;
  archiveEnvironmentV2(
    requestMessage: proto_environment_service_pb.ArchiveEnvironmentV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ArchiveEnvironmentV2Response|null) => void
  ): UnaryResponse;
  archiveEnvironmentV2(
    requestMessage: proto_environment_service_pb.ArchiveEnvironmentV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ArchiveEnvironmentV2Response|null) => void
  ): UnaryResponse;
  unarchiveEnvironmentV2(
    requestMessage: proto_environment_service_pb.UnarchiveEnvironmentV2Request,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UnarchiveEnvironmentV2Response|null) => void
  ): UnaryResponse;
  unarchiveEnvironmentV2(
    requestMessage: proto_environment_service_pb.UnarchiveEnvironmentV2Request,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UnarchiveEnvironmentV2Response|null) => void
  ): UnaryResponse;
  getProject(
    requestMessage: proto_environment_service_pb.GetProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.GetProjectResponse|null) => void
  ): UnaryResponse;
  getProject(
    requestMessage: proto_environment_service_pb.GetProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.GetProjectResponse|null) => void
  ): UnaryResponse;
  listProjects(
    requestMessage: proto_environment_service_pb.ListProjectsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ListProjectsResponse|null) => void
  ): UnaryResponse;
  listProjects(
    requestMessage: proto_environment_service_pb.ListProjectsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ListProjectsResponse|null) => void
  ): UnaryResponse;
  createProject(
    requestMessage: proto_environment_service_pb.CreateProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateProjectResponse|null) => void
  ): UnaryResponse;
  createProject(
    requestMessage: proto_environment_service_pb.CreateProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateProjectResponse|null) => void
  ): UnaryResponse;
  createTrialProject(
    requestMessage: proto_environment_service_pb.CreateTrialProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateTrialProjectResponse|null) => void
  ): UnaryResponse;
  createTrialProject(
    requestMessage: proto_environment_service_pb.CreateTrialProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateTrialProjectResponse|null) => void
  ): UnaryResponse;
  updateProject(
    requestMessage: proto_environment_service_pb.UpdateProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UpdateProjectResponse|null) => void
  ): UnaryResponse;
  updateProject(
    requestMessage: proto_environment_service_pb.UpdateProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UpdateProjectResponse|null) => void
  ): UnaryResponse;
  enableProject(
    requestMessage: proto_environment_service_pb.EnableProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.EnableProjectResponse|null) => void
  ): UnaryResponse;
  enableProject(
    requestMessage: proto_environment_service_pb.EnableProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.EnableProjectResponse|null) => void
  ): UnaryResponse;
  disableProject(
    requestMessage: proto_environment_service_pb.DisableProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.DisableProjectResponse|null) => void
  ): UnaryResponse;
  disableProject(
    requestMessage: proto_environment_service_pb.DisableProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.DisableProjectResponse|null) => void
  ): UnaryResponse;
  convertTrialProject(
    requestMessage: proto_environment_service_pb.ConvertTrialProjectRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ConvertTrialProjectResponse|null) => void
  ): UnaryResponse;
  convertTrialProject(
    requestMessage: proto_environment_service_pb.ConvertTrialProjectRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ConvertTrialProjectResponse|null) => void
  ): UnaryResponse;
  getOrganization(
    requestMessage: proto_environment_service_pb.GetOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.GetOrganizationResponse|null) => void
  ): UnaryResponse;
  getOrganization(
    requestMessage: proto_environment_service_pb.GetOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.GetOrganizationResponse|null) => void
  ): UnaryResponse;
  listOrganizations(
    requestMessage: proto_environment_service_pb.ListOrganizationsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ListOrganizationsResponse|null) => void
  ): UnaryResponse;
  listOrganizations(
    requestMessage: proto_environment_service_pb.ListOrganizationsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ListOrganizationsResponse|null) => void
  ): UnaryResponse;
  createOrganization(
    requestMessage: proto_environment_service_pb.CreateOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateOrganizationResponse|null) => void
  ): UnaryResponse;
  createOrganization(
    requestMessage: proto_environment_service_pb.CreateOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.CreateOrganizationResponse|null) => void
  ): UnaryResponse;
  updateOrganization(
    requestMessage: proto_environment_service_pb.UpdateOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UpdateOrganizationResponse|null) => void
  ): UnaryResponse;
  updateOrganization(
    requestMessage: proto_environment_service_pb.UpdateOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UpdateOrganizationResponse|null) => void
  ): UnaryResponse;
  enableOrganization(
    requestMessage: proto_environment_service_pb.EnableOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.EnableOrganizationResponse|null) => void
  ): UnaryResponse;
  enableOrganization(
    requestMessage: proto_environment_service_pb.EnableOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.EnableOrganizationResponse|null) => void
  ): UnaryResponse;
  disableOrganization(
    requestMessage: proto_environment_service_pb.DisableOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.DisableOrganizationResponse|null) => void
  ): UnaryResponse;
  disableOrganization(
    requestMessage: proto_environment_service_pb.DisableOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.DisableOrganizationResponse|null) => void
  ): UnaryResponse;
  archiveOrganization(
    requestMessage: proto_environment_service_pb.ArchiveOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ArchiveOrganizationResponse|null) => void
  ): UnaryResponse;
  archiveOrganization(
    requestMessage: proto_environment_service_pb.ArchiveOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ArchiveOrganizationResponse|null) => void
  ): UnaryResponse;
  unarchiveOrganization(
    requestMessage: proto_environment_service_pb.UnarchiveOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UnarchiveOrganizationResponse|null) => void
  ): UnaryResponse;
  unarchiveOrganization(
    requestMessage: proto_environment_service_pb.UnarchiveOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.UnarchiveOrganizationResponse|null) => void
  ): UnaryResponse;
  convertTrialOrganization(
    requestMessage: proto_environment_service_pb.ConvertTrialOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ConvertTrialOrganizationResponse|null) => void
  ): UnaryResponse;
  convertTrialOrganization(
    requestMessage: proto_environment_service_pb.ConvertTrialOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: proto_environment_service_pb.ConvertTrialOrganizationResponse|null) => void
  ): UnaryResponse;
}

