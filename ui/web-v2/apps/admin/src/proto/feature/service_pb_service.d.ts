// package: bucketeer.feature
// file: proto/feature/service.proto

import * as proto_feature_service_pb from "../../proto/feature/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type FeatureServiceGetFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.GetFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.GetFeatureResponse;
};

type FeatureServiceGetFeatures = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.GetFeaturesRequest;
  readonly responseType: typeof proto_feature_service_pb.GetFeaturesResponse;
};

type FeatureServiceListFeatures = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ListFeaturesRequest;
  readonly responseType: typeof proto_feature_service_pb.ListFeaturesResponse;
};

type FeatureServiceListEnabledFeatures = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ListEnabledFeaturesRequest;
  readonly responseType: typeof proto_feature_service_pb.ListEnabledFeaturesResponse;
};

type FeatureServiceCreateFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.CreateFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.CreateFeatureResponse;
};

type FeatureServiceEnableFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.EnableFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.EnableFeatureResponse;
};

type FeatureServiceDisableFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.DisableFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.DisableFeatureResponse;
};

type FeatureServiceArchiveFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ArchiveFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.ArchiveFeatureResponse;
};

type FeatureServiceUnarchiveFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.UnarchiveFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.UnarchiveFeatureResponse;
};

type FeatureServiceDeleteFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.DeleteFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.DeleteFeatureResponse;
};

type FeatureServiceUpdateFeatureDetails = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.UpdateFeatureDetailsRequest;
  readonly responseType: typeof proto_feature_service_pb.UpdateFeatureDetailsResponse;
};

type FeatureServiceUpdateFeatureVariations = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.UpdateFeatureVariationsRequest;
  readonly responseType: typeof proto_feature_service_pb.UpdateFeatureVariationsResponse;
};

type FeatureServiceUpdateFeatureTargeting = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.UpdateFeatureTargetingRequest;
  readonly responseType: typeof proto_feature_service_pb.UpdateFeatureTargetingResponse;
};

type FeatureServiceCloneFeature = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.CloneFeatureRequest;
  readonly responseType: typeof proto_feature_service_pb.CloneFeatureResponse;
};

type FeatureServiceCreateSegment = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.CreateSegmentRequest;
  readonly responseType: typeof proto_feature_service_pb.CreateSegmentResponse;
};

type FeatureServiceGetSegment = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.GetSegmentRequest;
  readonly responseType: typeof proto_feature_service_pb.GetSegmentResponse;
};

type FeatureServiceListSegments = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ListSegmentsRequest;
  readonly responseType: typeof proto_feature_service_pb.ListSegmentsResponse;
};

type FeatureServiceDeleteSegment = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.DeleteSegmentRequest;
  readonly responseType: typeof proto_feature_service_pb.DeleteSegmentResponse;
};

type FeatureServiceUpdateSegment = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.UpdateSegmentRequest;
  readonly responseType: typeof proto_feature_service_pb.UpdateSegmentResponse;
};

type FeatureServiceAddSegmentUser = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.AddSegmentUserRequest;
  readonly responseType: typeof proto_feature_service_pb.AddSegmentUserResponse;
};

type FeatureServiceDeleteSegmentUser = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.DeleteSegmentUserRequest;
  readonly responseType: typeof proto_feature_service_pb.DeleteSegmentUserResponse;
};

type FeatureServiceGetSegmentUser = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.GetSegmentUserRequest;
  readonly responseType: typeof proto_feature_service_pb.GetSegmentUserResponse;
};

type FeatureServiceListSegmentUsers = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ListSegmentUsersRequest;
  readonly responseType: typeof proto_feature_service_pb.ListSegmentUsersResponse;
};

type FeatureServiceBulkUploadSegmentUsers = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.BulkUploadSegmentUsersRequest;
  readonly responseType: typeof proto_feature_service_pb.BulkUploadSegmentUsersResponse;
};

type FeatureServiceBulkDownloadSegmentUsers = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.BulkDownloadSegmentUsersRequest;
  readonly responseType: typeof proto_feature_service_pb.BulkDownloadSegmentUsersResponse;
};

type FeatureServiceEvaluateFeatures = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.EvaluateFeaturesRequest;
  readonly responseType: typeof proto_feature_service_pb.EvaluateFeaturesResponse;
};

type FeatureServiceListTags = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ListTagsRequest;
  readonly responseType: typeof proto_feature_service_pb.ListTagsResponse;
};

type FeatureServiceCreateFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.CreateFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.CreateFlagTriggerResponse;
};

type FeatureServiceUpdateFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.UpdateFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.UpdateFlagTriggerResponse;
};

type FeatureServiceEnableFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.EnableFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.EnableFlagTriggerResponse;
};

type FeatureServiceDisableFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.DisableFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.DisableFlagTriggerResponse;
};

type FeatureServiceResetFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ResetFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.ResetFlagTriggerResponse;
};

type FeatureServiceDeleteFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.DeleteFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.DeleteFlagTriggerResponse;
};

type FeatureServiceGetFlagTrigger = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.GetFlagTriggerRequest;
  readonly responseType: typeof proto_feature_service_pb.GetFlagTriggerResponse;
};

type FeatureServiceListFlagTriggers = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.ListFlagTriggersRequest;
  readonly responseType: typeof proto_feature_service_pb.ListFlagTriggersResponse;
};

type FeatureServiceFlagTriggerWebhook = {
  readonly methodName: string;
  readonly service: typeof FeatureService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_feature_service_pb.FlagTriggerWebhookRequest;
  readonly responseType: typeof proto_feature_service_pb.FlagTriggerWebhookResponse;
};

export class FeatureService {
  static readonly serviceName: string;
  static readonly GetFeature: FeatureServiceGetFeature;
  static readonly GetFeatures: FeatureServiceGetFeatures;
  static readonly ListFeatures: FeatureServiceListFeatures;
  static readonly ListEnabledFeatures: FeatureServiceListEnabledFeatures;
  static readonly CreateFeature: FeatureServiceCreateFeature;
  static readonly EnableFeature: FeatureServiceEnableFeature;
  static readonly DisableFeature: FeatureServiceDisableFeature;
  static readonly ArchiveFeature: FeatureServiceArchiveFeature;
  static readonly UnarchiveFeature: FeatureServiceUnarchiveFeature;
  static readonly DeleteFeature: FeatureServiceDeleteFeature;
  static readonly UpdateFeatureDetails: FeatureServiceUpdateFeatureDetails;
  static readonly UpdateFeatureVariations: FeatureServiceUpdateFeatureVariations;
  static readonly UpdateFeatureTargeting: FeatureServiceUpdateFeatureTargeting;
  static readonly CloneFeature: FeatureServiceCloneFeature;
  static readonly CreateSegment: FeatureServiceCreateSegment;
  static readonly GetSegment: FeatureServiceGetSegment;
  static readonly ListSegments: FeatureServiceListSegments;
  static readonly DeleteSegment: FeatureServiceDeleteSegment;
  static readonly UpdateSegment: FeatureServiceUpdateSegment;
  static readonly AddSegmentUser: FeatureServiceAddSegmentUser;
  static readonly DeleteSegmentUser: FeatureServiceDeleteSegmentUser;
  static readonly GetSegmentUser: FeatureServiceGetSegmentUser;
  static readonly ListSegmentUsers: FeatureServiceListSegmentUsers;
  static readonly BulkUploadSegmentUsers: FeatureServiceBulkUploadSegmentUsers;
  static readonly BulkDownloadSegmentUsers: FeatureServiceBulkDownloadSegmentUsers;
  static readonly EvaluateFeatures: FeatureServiceEvaluateFeatures;
  static readonly ListTags: FeatureServiceListTags;
  static readonly CreateFlagTrigger: FeatureServiceCreateFlagTrigger;
  static readonly UpdateFlagTrigger: FeatureServiceUpdateFlagTrigger;
  static readonly EnableFlagTrigger: FeatureServiceEnableFlagTrigger;
  static readonly DisableFlagTrigger: FeatureServiceDisableFlagTrigger;
  static readonly ResetFlagTrigger: FeatureServiceResetFlagTrigger;
  static readonly DeleteFlagTrigger: FeatureServiceDeleteFlagTrigger;
  static readonly GetFlagTrigger: FeatureServiceGetFlagTrigger;
  static readonly ListFlagTriggers: FeatureServiceListFlagTriggers;
  static readonly FlagTriggerWebhook: FeatureServiceFlagTriggerWebhook;
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

export class FeatureServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getFeature(
    requestMessage: proto_feature_service_pb.GetFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetFeatureResponse|null) => void
  ): UnaryResponse;
  getFeature(
    requestMessage: proto_feature_service_pb.GetFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetFeatureResponse|null) => void
  ): UnaryResponse;
  getFeatures(
    requestMessage: proto_feature_service_pb.GetFeaturesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetFeaturesResponse|null) => void
  ): UnaryResponse;
  getFeatures(
    requestMessage: proto_feature_service_pb.GetFeaturesRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetFeaturesResponse|null) => void
  ): UnaryResponse;
  listFeatures(
    requestMessage: proto_feature_service_pb.ListFeaturesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListFeaturesResponse|null) => void
  ): UnaryResponse;
  listFeatures(
    requestMessage: proto_feature_service_pb.ListFeaturesRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListFeaturesResponse|null) => void
  ): UnaryResponse;
  listEnabledFeatures(
    requestMessage: proto_feature_service_pb.ListEnabledFeaturesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListEnabledFeaturesResponse|null) => void
  ): UnaryResponse;
  listEnabledFeatures(
    requestMessage: proto_feature_service_pb.ListEnabledFeaturesRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListEnabledFeaturesResponse|null) => void
  ): UnaryResponse;
  createFeature(
    requestMessage: proto_feature_service_pb.CreateFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CreateFeatureResponse|null) => void
  ): UnaryResponse;
  createFeature(
    requestMessage: proto_feature_service_pb.CreateFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CreateFeatureResponse|null) => void
  ): UnaryResponse;
  enableFeature(
    requestMessage: proto_feature_service_pb.EnableFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.EnableFeatureResponse|null) => void
  ): UnaryResponse;
  enableFeature(
    requestMessage: proto_feature_service_pb.EnableFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.EnableFeatureResponse|null) => void
  ): UnaryResponse;
  disableFeature(
    requestMessage: proto_feature_service_pb.DisableFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DisableFeatureResponse|null) => void
  ): UnaryResponse;
  disableFeature(
    requestMessage: proto_feature_service_pb.DisableFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DisableFeatureResponse|null) => void
  ): UnaryResponse;
  archiveFeature(
    requestMessage: proto_feature_service_pb.ArchiveFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ArchiveFeatureResponse|null) => void
  ): UnaryResponse;
  archiveFeature(
    requestMessage: proto_feature_service_pb.ArchiveFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ArchiveFeatureResponse|null) => void
  ): UnaryResponse;
  unarchiveFeature(
    requestMessage: proto_feature_service_pb.UnarchiveFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UnarchiveFeatureResponse|null) => void
  ): UnaryResponse;
  unarchiveFeature(
    requestMessage: proto_feature_service_pb.UnarchiveFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UnarchiveFeatureResponse|null) => void
  ): UnaryResponse;
  deleteFeature(
    requestMessage: proto_feature_service_pb.DeleteFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteFeatureResponse|null) => void
  ): UnaryResponse;
  deleteFeature(
    requestMessage: proto_feature_service_pb.DeleteFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteFeatureResponse|null) => void
  ): UnaryResponse;
  updateFeatureDetails(
    requestMessage: proto_feature_service_pb.UpdateFeatureDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFeatureDetailsResponse|null) => void
  ): UnaryResponse;
  updateFeatureDetails(
    requestMessage: proto_feature_service_pb.UpdateFeatureDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFeatureDetailsResponse|null) => void
  ): UnaryResponse;
  updateFeatureVariations(
    requestMessage: proto_feature_service_pb.UpdateFeatureVariationsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFeatureVariationsResponse|null) => void
  ): UnaryResponse;
  updateFeatureVariations(
    requestMessage: proto_feature_service_pb.UpdateFeatureVariationsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFeatureVariationsResponse|null) => void
  ): UnaryResponse;
  updateFeatureTargeting(
    requestMessage: proto_feature_service_pb.UpdateFeatureTargetingRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFeatureTargetingResponse|null) => void
  ): UnaryResponse;
  updateFeatureTargeting(
    requestMessage: proto_feature_service_pb.UpdateFeatureTargetingRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFeatureTargetingResponse|null) => void
  ): UnaryResponse;
  cloneFeature(
    requestMessage: proto_feature_service_pb.CloneFeatureRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CloneFeatureResponse|null) => void
  ): UnaryResponse;
  cloneFeature(
    requestMessage: proto_feature_service_pb.CloneFeatureRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CloneFeatureResponse|null) => void
  ): UnaryResponse;
  createSegment(
    requestMessage: proto_feature_service_pb.CreateSegmentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CreateSegmentResponse|null) => void
  ): UnaryResponse;
  createSegment(
    requestMessage: proto_feature_service_pb.CreateSegmentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CreateSegmentResponse|null) => void
  ): UnaryResponse;
  getSegment(
    requestMessage: proto_feature_service_pb.GetSegmentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetSegmentResponse|null) => void
  ): UnaryResponse;
  getSegment(
    requestMessage: proto_feature_service_pb.GetSegmentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetSegmentResponse|null) => void
  ): UnaryResponse;
  listSegments(
    requestMessage: proto_feature_service_pb.ListSegmentsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListSegmentsResponse|null) => void
  ): UnaryResponse;
  listSegments(
    requestMessage: proto_feature_service_pb.ListSegmentsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListSegmentsResponse|null) => void
  ): UnaryResponse;
  deleteSegment(
    requestMessage: proto_feature_service_pb.DeleteSegmentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteSegmentResponse|null) => void
  ): UnaryResponse;
  deleteSegment(
    requestMessage: proto_feature_service_pb.DeleteSegmentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteSegmentResponse|null) => void
  ): UnaryResponse;
  updateSegment(
    requestMessage: proto_feature_service_pb.UpdateSegmentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateSegmentResponse|null) => void
  ): UnaryResponse;
  updateSegment(
    requestMessage: proto_feature_service_pb.UpdateSegmentRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateSegmentResponse|null) => void
  ): UnaryResponse;
  addSegmentUser(
    requestMessage: proto_feature_service_pb.AddSegmentUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.AddSegmentUserResponse|null) => void
  ): UnaryResponse;
  addSegmentUser(
    requestMessage: proto_feature_service_pb.AddSegmentUserRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.AddSegmentUserResponse|null) => void
  ): UnaryResponse;
  deleteSegmentUser(
    requestMessage: proto_feature_service_pb.DeleteSegmentUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteSegmentUserResponse|null) => void
  ): UnaryResponse;
  deleteSegmentUser(
    requestMessage: proto_feature_service_pb.DeleteSegmentUserRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteSegmentUserResponse|null) => void
  ): UnaryResponse;
  getSegmentUser(
    requestMessage: proto_feature_service_pb.GetSegmentUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetSegmentUserResponse|null) => void
  ): UnaryResponse;
  getSegmentUser(
    requestMessage: proto_feature_service_pb.GetSegmentUserRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetSegmentUserResponse|null) => void
  ): UnaryResponse;
  listSegmentUsers(
    requestMessage: proto_feature_service_pb.ListSegmentUsersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListSegmentUsersResponse|null) => void
  ): UnaryResponse;
  listSegmentUsers(
    requestMessage: proto_feature_service_pb.ListSegmentUsersRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListSegmentUsersResponse|null) => void
  ): UnaryResponse;
  bulkUploadSegmentUsers(
    requestMessage: proto_feature_service_pb.BulkUploadSegmentUsersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.BulkUploadSegmentUsersResponse|null) => void
  ): UnaryResponse;
  bulkUploadSegmentUsers(
    requestMessage: proto_feature_service_pb.BulkUploadSegmentUsersRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.BulkUploadSegmentUsersResponse|null) => void
  ): UnaryResponse;
  bulkDownloadSegmentUsers(
    requestMessage: proto_feature_service_pb.BulkDownloadSegmentUsersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.BulkDownloadSegmentUsersResponse|null) => void
  ): UnaryResponse;
  bulkDownloadSegmentUsers(
    requestMessage: proto_feature_service_pb.BulkDownloadSegmentUsersRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.BulkDownloadSegmentUsersResponse|null) => void
  ): UnaryResponse;
  evaluateFeatures(
    requestMessage: proto_feature_service_pb.EvaluateFeaturesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.EvaluateFeaturesResponse|null) => void
  ): UnaryResponse;
  evaluateFeatures(
    requestMessage: proto_feature_service_pb.EvaluateFeaturesRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.EvaluateFeaturesResponse|null) => void
  ): UnaryResponse;
  listTags(
    requestMessage: proto_feature_service_pb.ListTagsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListTagsResponse|null) => void
  ): UnaryResponse;
  listTags(
    requestMessage: proto_feature_service_pb.ListTagsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListTagsResponse|null) => void
  ): UnaryResponse;
  createFlagTrigger(
    requestMessage: proto_feature_service_pb.CreateFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CreateFlagTriggerResponse|null) => void
  ): UnaryResponse;
  createFlagTrigger(
    requestMessage: proto_feature_service_pb.CreateFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.CreateFlagTriggerResponse|null) => void
  ): UnaryResponse;
  updateFlagTrigger(
    requestMessage: proto_feature_service_pb.UpdateFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFlagTriggerResponse|null) => void
  ): UnaryResponse;
  updateFlagTrigger(
    requestMessage: proto_feature_service_pb.UpdateFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.UpdateFlagTriggerResponse|null) => void
  ): UnaryResponse;
  enableFlagTrigger(
    requestMessage: proto_feature_service_pb.EnableFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.EnableFlagTriggerResponse|null) => void
  ): UnaryResponse;
  enableFlagTrigger(
    requestMessage: proto_feature_service_pb.EnableFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.EnableFlagTriggerResponse|null) => void
  ): UnaryResponse;
  disableFlagTrigger(
    requestMessage: proto_feature_service_pb.DisableFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DisableFlagTriggerResponse|null) => void
  ): UnaryResponse;
  disableFlagTrigger(
    requestMessage: proto_feature_service_pb.DisableFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DisableFlagTriggerResponse|null) => void
  ): UnaryResponse;
  resetFlagTrigger(
    requestMessage: proto_feature_service_pb.ResetFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ResetFlagTriggerResponse|null) => void
  ): UnaryResponse;
  resetFlagTrigger(
    requestMessage: proto_feature_service_pb.ResetFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ResetFlagTriggerResponse|null) => void
  ): UnaryResponse;
  deleteFlagTrigger(
    requestMessage: proto_feature_service_pb.DeleteFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteFlagTriggerResponse|null) => void
  ): UnaryResponse;
  deleteFlagTrigger(
    requestMessage: proto_feature_service_pb.DeleteFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.DeleteFlagTriggerResponse|null) => void
  ): UnaryResponse;
  getFlagTrigger(
    requestMessage: proto_feature_service_pb.GetFlagTriggerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetFlagTriggerResponse|null) => void
  ): UnaryResponse;
  getFlagTrigger(
    requestMessage: proto_feature_service_pb.GetFlagTriggerRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.GetFlagTriggerResponse|null) => void
  ): UnaryResponse;
  listFlagTriggers(
    requestMessage: proto_feature_service_pb.ListFlagTriggersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListFlagTriggersResponse|null) => void
  ): UnaryResponse;
  listFlagTriggers(
    requestMessage: proto_feature_service_pb.ListFlagTriggersRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.ListFlagTriggersResponse|null) => void
  ): UnaryResponse;
  flagTriggerWebhook(
    requestMessage: proto_feature_service_pb.FlagTriggerWebhookRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.FlagTriggerWebhookResponse|null) => void
  ): UnaryResponse;
  flagTriggerWebhook(
    requestMessage: proto_feature_service_pb.FlagTriggerWebhookRequest,
    callback: (error: ServiceError|null, responseMessage: proto_feature_service_pb.FlagTriggerWebhookResponse|null) => void
  ): UnaryResponse;
}

