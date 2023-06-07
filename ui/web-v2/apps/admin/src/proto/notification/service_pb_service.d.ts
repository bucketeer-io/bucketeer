// package: bucketeer.notification
// file: proto/notification/service.proto

import * as proto_notification_service_pb from "../../proto/notification/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type NotificationServiceGetAdminSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.GetAdminSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.GetAdminSubscriptionResponse;
};

type NotificationServiceListAdminSubscriptions = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.ListAdminSubscriptionsRequest;
  readonly responseType: typeof proto_notification_service_pb.ListAdminSubscriptionsResponse;
};

type NotificationServiceListEnabledAdminSubscriptions = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.ListEnabledAdminSubscriptionsRequest;
  readonly responseType: typeof proto_notification_service_pb.ListEnabledAdminSubscriptionsResponse;
};

type NotificationServiceCreateAdminSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.CreateAdminSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.CreateAdminSubscriptionResponse;
};

type NotificationServiceDeleteAdminSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.DeleteAdminSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.DeleteAdminSubscriptionResponse;
};

type NotificationServiceEnableAdminSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.EnableAdminSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.EnableAdminSubscriptionResponse;
};

type NotificationServiceDisableAdminSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.DisableAdminSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.DisableAdminSubscriptionResponse;
};

type NotificationServiceUpdateAdminSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.UpdateAdminSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.UpdateAdminSubscriptionResponse;
};

type NotificationServiceGetSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.GetSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.GetSubscriptionResponse;
};

type NotificationServiceListSubscriptions = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.ListSubscriptionsRequest;
  readonly responseType: typeof proto_notification_service_pb.ListSubscriptionsResponse;
};

type NotificationServiceListEnabledSubscriptions = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.ListEnabledSubscriptionsRequest;
  readonly responseType: typeof proto_notification_service_pb.ListEnabledSubscriptionsResponse;
};

type NotificationServiceCreateSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.CreateSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.CreateSubscriptionResponse;
};

type NotificationServiceDeleteSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.DeleteSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.DeleteSubscriptionResponse;
};

type NotificationServiceEnableSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.EnableSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.EnableSubscriptionResponse;
};

type NotificationServiceDisableSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.DisableSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.DisableSubscriptionResponse;
};

type NotificationServiceUpdateSubscription = {
  readonly methodName: string;
  readonly service: typeof NotificationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_notification_service_pb.UpdateSubscriptionRequest;
  readonly responseType: typeof proto_notification_service_pb.UpdateSubscriptionResponse;
};

export class NotificationService {
  static readonly serviceName: string;
  static readonly GetAdminSubscription: NotificationServiceGetAdminSubscription;
  static readonly ListAdminSubscriptions: NotificationServiceListAdminSubscriptions;
  static readonly ListEnabledAdminSubscriptions: NotificationServiceListEnabledAdminSubscriptions;
  static readonly CreateAdminSubscription: NotificationServiceCreateAdminSubscription;
  static readonly DeleteAdminSubscription: NotificationServiceDeleteAdminSubscription;
  static readonly EnableAdminSubscription: NotificationServiceEnableAdminSubscription;
  static readonly DisableAdminSubscription: NotificationServiceDisableAdminSubscription;
  static readonly UpdateAdminSubscription: NotificationServiceUpdateAdminSubscription;
  static readonly GetSubscription: NotificationServiceGetSubscription;
  static readonly ListSubscriptions: NotificationServiceListSubscriptions;
  static readonly ListEnabledSubscriptions: NotificationServiceListEnabledSubscriptions;
  static readonly CreateSubscription: NotificationServiceCreateSubscription;
  static readonly DeleteSubscription: NotificationServiceDeleteSubscription;
  static readonly EnableSubscription: NotificationServiceEnableSubscription;
  static readonly DisableSubscription: NotificationServiceDisableSubscription;
  static readonly UpdateSubscription: NotificationServiceUpdateSubscription;
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

export class NotificationServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getAdminSubscription(
    requestMessage: proto_notification_service_pb.GetAdminSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.GetAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  getAdminSubscription(
    requestMessage: proto_notification_service_pb.GetAdminSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.GetAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  listAdminSubscriptions(
    requestMessage: proto_notification_service_pb.ListAdminSubscriptionsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListAdminSubscriptionsResponse|null) => void
  ): UnaryResponse;
  listAdminSubscriptions(
    requestMessage: proto_notification_service_pb.ListAdminSubscriptionsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListAdminSubscriptionsResponse|null) => void
  ): UnaryResponse;
  listEnabledAdminSubscriptions(
    requestMessage: proto_notification_service_pb.ListEnabledAdminSubscriptionsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListEnabledAdminSubscriptionsResponse|null) => void
  ): UnaryResponse;
  listEnabledAdminSubscriptions(
    requestMessage: proto_notification_service_pb.ListEnabledAdminSubscriptionsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListEnabledAdminSubscriptionsResponse|null) => void
  ): UnaryResponse;
  createAdminSubscription(
    requestMessage: proto_notification_service_pb.CreateAdminSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.CreateAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  createAdminSubscription(
    requestMessage: proto_notification_service_pb.CreateAdminSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.CreateAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  deleteAdminSubscription(
    requestMessage: proto_notification_service_pb.DeleteAdminSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DeleteAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  deleteAdminSubscription(
    requestMessage: proto_notification_service_pb.DeleteAdminSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DeleteAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  enableAdminSubscription(
    requestMessage: proto_notification_service_pb.EnableAdminSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.EnableAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  enableAdminSubscription(
    requestMessage: proto_notification_service_pb.EnableAdminSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.EnableAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  disableAdminSubscription(
    requestMessage: proto_notification_service_pb.DisableAdminSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DisableAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  disableAdminSubscription(
    requestMessage: proto_notification_service_pb.DisableAdminSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DisableAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  updateAdminSubscription(
    requestMessage: proto_notification_service_pb.UpdateAdminSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.UpdateAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  updateAdminSubscription(
    requestMessage: proto_notification_service_pb.UpdateAdminSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.UpdateAdminSubscriptionResponse|null) => void
  ): UnaryResponse;
  getSubscription(
    requestMessage: proto_notification_service_pb.GetSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.GetSubscriptionResponse|null) => void
  ): UnaryResponse;
  getSubscription(
    requestMessage: proto_notification_service_pb.GetSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.GetSubscriptionResponse|null) => void
  ): UnaryResponse;
  listSubscriptions(
    requestMessage: proto_notification_service_pb.ListSubscriptionsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListSubscriptionsResponse|null) => void
  ): UnaryResponse;
  listSubscriptions(
    requestMessage: proto_notification_service_pb.ListSubscriptionsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListSubscriptionsResponse|null) => void
  ): UnaryResponse;
  listEnabledSubscriptions(
    requestMessage: proto_notification_service_pb.ListEnabledSubscriptionsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListEnabledSubscriptionsResponse|null) => void
  ): UnaryResponse;
  listEnabledSubscriptions(
    requestMessage: proto_notification_service_pb.ListEnabledSubscriptionsRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.ListEnabledSubscriptionsResponse|null) => void
  ): UnaryResponse;
  createSubscription(
    requestMessage: proto_notification_service_pb.CreateSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.CreateSubscriptionResponse|null) => void
  ): UnaryResponse;
  createSubscription(
    requestMessage: proto_notification_service_pb.CreateSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.CreateSubscriptionResponse|null) => void
  ): UnaryResponse;
  deleteSubscription(
    requestMessage: proto_notification_service_pb.DeleteSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DeleteSubscriptionResponse|null) => void
  ): UnaryResponse;
  deleteSubscription(
    requestMessage: proto_notification_service_pb.DeleteSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DeleteSubscriptionResponse|null) => void
  ): UnaryResponse;
  enableSubscription(
    requestMessage: proto_notification_service_pb.EnableSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.EnableSubscriptionResponse|null) => void
  ): UnaryResponse;
  enableSubscription(
    requestMessage: proto_notification_service_pb.EnableSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.EnableSubscriptionResponse|null) => void
  ): UnaryResponse;
  disableSubscription(
    requestMessage: proto_notification_service_pb.DisableSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DisableSubscriptionResponse|null) => void
  ): UnaryResponse;
  disableSubscription(
    requestMessage: proto_notification_service_pb.DisableSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.DisableSubscriptionResponse|null) => void
  ): UnaryResponse;
  updateSubscription(
    requestMessage: proto_notification_service_pb.UpdateSubscriptionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.UpdateSubscriptionResponse|null) => void
  ): UnaryResponse;
  updateSubscription(
    requestMessage: proto_notification_service_pb.UpdateSubscriptionRequest,
    callback: (error: ServiceError|null, responseMessage: proto_notification_service_pb.UpdateSubscriptionResponse|null) => void
  ): UnaryResponse;
}

