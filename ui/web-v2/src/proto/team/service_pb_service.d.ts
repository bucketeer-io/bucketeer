// package: bucketeer.team
// file: proto/team/service.proto

import * as proto_team_service_pb from '../../proto/team/service_pb';
import { grpc } from '@improbable-eng/grpc-web';

type TeamServiceCreateTeam = {
  readonly methodName: string;
  readonly service: typeof TeamService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_team_service_pb.CreateTeamRequest;
  readonly responseType: typeof proto_team_service_pb.CreateTeamResponse;
};

type TeamServiceDeleteTeam = {
  readonly methodName: string;
  readonly service: typeof TeamService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_team_service_pb.DeleteTeamRequest;
  readonly responseType: typeof proto_team_service_pb.DeleteTeamResponse;
};

type TeamServiceListTeams = {
  readonly methodName: string;
  readonly service: typeof TeamService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_team_service_pb.ListTeamsRequest;
  readonly responseType: typeof proto_team_service_pb.ListTeamsResponse;
};

export class TeamService {
  static readonly serviceName: string;
  static readonly CreateTeam: TeamServiceCreateTeam;
  static readonly DeleteTeam: TeamServiceDeleteTeam;
  static readonly ListTeams: TeamServiceListTeams;
}

export type ServiceError = {
  message: string;
  code: number;
  metadata: grpc.Metadata;
};
export type Status = { details: string; code: number; metadata: grpc.Metadata };

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
  on(
    type: 'data',
    handler: (message: ResT) => void
  ): BidirectionalStream<ReqT, ResT>;
  on(
    type: 'end',
    handler: (status?: Status) => void
  ): BidirectionalStream<ReqT, ResT>;
  on(
    type: 'status',
    handler: (status: Status) => void
  ): BidirectionalStream<ReqT, ResT>;
}

export class TeamServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createTeam(
    requestMessage: proto_team_service_pb.CreateTeamRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_team_service_pb.CreateTeamResponse | null
    ) => void
  ): UnaryResponse;
  createTeam(
    requestMessage: proto_team_service_pb.CreateTeamRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_team_service_pb.CreateTeamResponse | null
    ) => void
  ): UnaryResponse;
  deleteTeam(
    requestMessage: proto_team_service_pb.DeleteTeamRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_team_service_pb.DeleteTeamResponse | null
    ) => void
  ): UnaryResponse;
  deleteTeam(
    requestMessage: proto_team_service_pb.DeleteTeamRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_team_service_pb.DeleteTeamResponse | null
    ) => void
  ): UnaryResponse;
  listTeams(
    requestMessage: proto_team_service_pb.ListTeamsRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_team_service_pb.ListTeamsResponse | null
    ) => void
  ): UnaryResponse;
  listTeams(
    requestMessage: proto_team_service_pb.ListTeamsRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_team_service_pb.ListTeamsResponse | null
    ) => void
  ): UnaryResponse;
}
