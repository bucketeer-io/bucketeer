// package: bucketeer.autoops
// file: proto/autoops/service.proto

import * as proto_autoops_service_pb from "../../proto/autoops/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type AutoOpsServiceGetAutoOpsRule = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.GetAutoOpsRuleRequest;
  readonly responseType: typeof proto_autoops_service_pb.GetAutoOpsRuleResponse;
};

type AutoOpsServiceListAutoOpsRules = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.ListAutoOpsRulesRequest;
  readonly responseType: typeof proto_autoops_service_pb.ListAutoOpsRulesResponse;
};

type AutoOpsServiceCreateAutoOpsRule = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.CreateAutoOpsRuleRequest;
  readonly responseType: typeof proto_autoops_service_pb.CreateAutoOpsRuleResponse;
};

type AutoOpsServiceDeleteAutoOpsRule = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.DeleteAutoOpsRuleRequest;
  readonly responseType: typeof proto_autoops_service_pb.DeleteAutoOpsRuleResponse;
};

type AutoOpsServiceUpdateAutoOpsRule = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.UpdateAutoOpsRuleRequest;
  readonly responseType: typeof proto_autoops_service_pb.UpdateAutoOpsRuleResponse;
};

type AutoOpsServiceExecuteAutoOps = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.ExecuteAutoOpsRequest;
  readonly responseType: typeof proto_autoops_service_pb.ExecuteAutoOpsResponse;
};

type AutoOpsServiceListOpsCounts = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.ListOpsCountsRequest;
  readonly responseType: typeof proto_autoops_service_pb.ListOpsCountsResponse;
};

type AutoOpsServiceCreateWebhook = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.CreateWebhookRequest;
  readonly responseType: typeof proto_autoops_service_pb.CreateWebhookResponse;
};

type AutoOpsServiceGetWebhook = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.GetWebhookRequest;
  readonly responseType: typeof proto_autoops_service_pb.GetWebhookResponse;
};

type AutoOpsServiceUpdateWebhook = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.UpdateWebhookRequest;
  readonly responseType: typeof proto_autoops_service_pb.UpdateWebhookResponse;
};

type AutoOpsServiceDeleteWebhook = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.DeleteWebhookRequest;
  readonly responseType: typeof proto_autoops_service_pb.DeleteWebhookResponse;
};

type AutoOpsServiceListWebhooks = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.ListWebhooksRequest;
  readonly responseType: typeof proto_autoops_service_pb.ListWebhooksResponse;
};

type AutoOpsServiceCreateProgressiveRollout = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.CreateProgressiveRolloutRequest;
  readonly responseType: typeof proto_autoops_service_pb.CreateProgressiveRolloutResponse;
};

type AutoOpsServiceGetProgressiveRollout = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.GetProgressiveRolloutRequest;
  readonly responseType: typeof proto_autoops_service_pb.GetProgressiveRolloutResponse;
};

type AutoOpsServiceDeleteProgressiveRollout = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.DeleteProgressiveRolloutRequest;
  readonly responseType: typeof proto_autoops_service_pb.DeleteProgressiveRolloutResponse;
};

type AutoOpsServiceListProgressiveRollouts = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.ListProgressiveRolloutsRequest;
  readonly responseType: typeof proto_autoops_service_pb.ListProgressiveRolloutsResponse;
};

type AutoOpsServiceExecuteProgressiveRollout = {
  readonly methodName: string;
  readonly service: typeof AutoOpsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_autoops_service_pb.ExecuteProgressiveRolloutRequest;
  readonly responseType: typeof proto_autoops_service_pb.ExecuteProgressiveRolloutResponse;
};

export class AutoOpsService {
  static readonly serviceName: string;
  static readonly GetAutoOpsRule: AutoOpsServiceGetAutoOpsRule;
  static readonly ListAutoOpsRules: AutoOpsServiceListAutoOpsRules;
  static readonly CreateAutoOpsRule: AutoOpsServiceCreateAutoOpsRule;
  static readonly DeleteAutoOpsRule: AutoOpsServiceDeleteAutoOpsRule;
  static readonly UpdateAutoOpsRule: AutoOpsServiceUpdateAutoOpsRule;
  static readonly ExecuteAutoOps: AutoOpsServiceExecuteAutoOps;
  static readonly ListOpsCounts: AutoOpsServiceListOpsCounts;
  static readonly CreateWebhook: AutoOpsServiceCreateWebhook;
  static readonly GetWebhook: AutoOpsServiceGetWebhook;
  static readonly UpdateWebhook: AutoOpsServiceUpdateWebhook;
  static readonly DeleteWebhook: AutoOpsServiceDeleteWebhook;
  static readonly ListWebhooks: AutoOpsServiceListWebhooks;
  static readonly CreateProgressiveRollout: AutoOpsServiceCreateProgressiveRollout;
  static readonly GetProgressiveRollout: AutoOpsServiceGetProgressiveRollout;
  static readonly DeleteProgressiveRollout: AutoOpsServiceDeleteProgressiveRollout;
  static readonly ListProgressiveRollouts: AutoOpsServiceListProgressiveRollouts;
  static readonly ExecuteProgressiveRollout: AutoOpsServiceExecuteProgressiveRollout;
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

export class AutoOpsServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getAutoOpsRule(
    requestMessage: proto_autoops_service_pb.GetAutoOpsRuleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.GetAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  getAutoOpsRule(
    requestMessage: proto_autoops_service_pb.GetAutoOpsRuleRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.GetAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  listAutoOpsRules(
    requestMessage: proto_autoops_service_pb.ListAutoOpsRulesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListAutoOpsRulesResponse|null) => void
  ): UnaryResponse;
  listAutoOpsRules(
    requestMessage: proto_autoops_service_pb.ListAutoOpsRulesRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListAutoOpsRulesResponse|null) => void
  ): UnaryResponse;
  createAutoOpsRule(
    requestMessage: proto_autoops_service_pb.CreateAutoOpsRuleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.CreateAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  createAutoOpsRule(
    requestMessage: proto_autoops_service_pb.CreateAutoOpsRuleRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.CreateAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  deleteAutoOpsRule(
    requestMessage: proto_autoops_service_pb.DeleteAutoOpsRuleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.DeleteAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  deleteAutoOpsRule(
    requestMessage: proto_autoops_service_pb.DeleteAutoOpsRuleRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.DeleteAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  updateAutoOpsRule(
    requestMessage: proto_autoops_service_pb.UpdateAutoOpsRuleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.UpdateAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  updateAutoOpsRule(
    requestMessage: proto_autoops_service_pb.UpdateAutoOpsRuleRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.UpdateAutoOpsRuleResponse|null) => void
  ): UnaryResponse;
  executeAutoOps(
    requestMessage: proto_autoops_service_pb.ExecuteAutoOpsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ExecuteAutoOpsResponse|null) => void
  ): UnaryResponse;
  executeAutoOps(
    requestMessage: proto_autoops_service_pb.ExecuteAutoOpsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ExecuteAutoOpsResponse|null) => void
  ): UnaryResponse;
  listOpsCounts(
    requestMessage: proto_autoops_service_pb.ListOpsCountsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListOpsCountsResponse|null) => void
  ): UnaryResponse;
  listOpsCounts(
    requestMessage: proto_autoops_service_pb.ListOpsCountsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListOpsCountsResponse|null) => void
  ): UnaryResponse;
  createWebhook(
    requestMessage: proto_autoops_service_pb.CreateWebhookRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.CreateWebhookResponse|null) => void
  ): UnaryResponse;
  createWebhook(
    requestMessage: proto_autoops_service_pb.CreateWebhookRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.CreateWebhookResponse|null) => void
  ): UnaryResponse;
  getWebhook(
    requestMessage: proto_autoops_service_pb.GetWebhookRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.GetWebhookResponse|null) => void
  ): UnaryResponse;
  getWebhook(
    requestMessage: proto_autoops_service_pb.GetWebhookRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.GetWebhookResponse|null) => void
  ): UnaryResponse;
  updateWebhook(
    requestMessage: proto_autoops_service_pb.UpdateWebhookRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.UpdateWebhookResponse|null) => void
  ): UnaryResponse;
  updateWebhook(
    requestMessage: proto_autoops_service_pb.UpdateWebhookRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.UpdateWebhookResponse|null) => void
  ): UnaryResponse;
  deleteWebhook(
    requestMessage: proto_autoops_service_pb.DeleteWebhookRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.DeleteWebhookResponse|null) => void
  ): UnaryResponse;
  deleteWebhook(
    requestMessage: proto_autoops_service_pb.DeleteWebhookRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.DeleteWebhookResponse|null) => void
  ): UnaryResponse;
  listWebhooks(
    requestMessage: proto_autoops_service_pb.ListWebhooksRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListWebhooksResponse|null) => void
  ): UnaryResponse;
  listWebhooks(
    requestMessage: proto_autoops_service_pb.ListWebhooksRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListWebhooksResponse|null) => void
  ): UnaryResponse;
  createProgressiveRollout(
    requestMessage: proto_autoops_service_pb.CreateProgressiveRolloutRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.CreateProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  createProgressiveRollout(
    requestMessage: proto_autoops_service_pb.CreateProgressiveRolloutRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.CreateProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  getProgressiveRollout(
    requestMessage: proto_autoops_service_pb.GetProgressiveRolloutRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.GetProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  getProgressiveRollout(
    requestMessage: proto_autoops_service_pb.GetProgressiveRolloutRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.GetProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  deleteProgressiveRollout(
    requestMessage: proto_autoops_service_pb.DeleteProgressiveRolloutRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.DeleteProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  deleteProgressiveRollout(
    requestMessage: proto_autoops_service_pb.DeleteProgressiveRolloutRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.DeleteProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  listProgressiveRollouts(
    requestMessage: proto_autoops_service_pb.ListProgressiveRolloutsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListProgressiveRolloutsResponse|null) => void
  ): UnaryResponse;
  listProgressiveRollouts(
    requestMessage: proto_autoops_service_pb.ListProgressiveRolloutsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ListProgressiveRolloutsResponse|null) => void
  ): UnaryResponse;
  executeProgressiveRollout(
    requestMessage: proto_autoops_service_pb.ExecuteProgressiveRolloutRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ExecuteProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
  executeProgressiveRollout(
    requestMessage: proto_autoops_service_pb.ExecuteProgressiveRolloutRequest,
    callback: (error: ServiceError|null, responseMessage: proto_autoops_service_pb.ExecuteProgressiveRolloutResponse|null) => void
  ): UnaryResponse;
}

