import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';
import { Undefinable } from 'option-t/lib/Undefinable/Undefinable';

import { urls } from '../config';
import {
  CreateAdminSubscriptionRequest,
  CreateAdminSubscriptionResponse,
  DeleteAdminSubscriptionRequest,
  DeleteAdminSubscriptionResponse,
  DisableAdminSubscriptionRequest,
  DisableAdminSubscriptionResponse,
  EnableAdminSubscriptionRequest,
  EnableAdminSubscriptionResponse,
  GetAdminSubscriptionRequest,
  GetAdminSubscriptionResponse,
  ListAdminSubscriptionsRequest,
  ListAdminSubscriptionsResponse,
  UpdateAdminSubscriptionRequest,
  UpdateAdminSubscriptionResponse,
} from '../proto/notification/service_pb';
import {
  NotificationServiceClient,
  ServiceError,
} from '../proto/notification/service_pb_service';

import { getMetaDataForClient as getMetaData } from './utils';

export class NotificationServiceError<Request> extends Error {
  request: Request;
  error: Nullable<ServiceError>;
  constructor(
    message: Undefinable<string>,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, NotificationServiceError);
    }
    this.name = 'NotificationServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new NotificationServiceClient(urls.GRPC);

export interface CreateAdminSubscriptionResult {
  request: CreateAdminSubscriptionRequest;
  response: CreateAdminSubscriptionResponse;
}

export function createSubscription(
  request: CreateAdminSubscriptionRequest
): Promise<CreateAdminSubscriptionResult> {
  return new Promise(
    (
      resolve: (result: CreateAdminSubscriptionResult) => void,
      reject
    ): void => {
      client.createAdminSubscription(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface ListAdminSubscriptionsResult {
  request: ListAdminSubscriptionsRequest;
  response: ListAdminSubscriptionsResponse;
}

export function listSubscriptions(
  request: ListAdminSubscriptionsRequest
): Promise<ListAdminSubscriptionsResult> {
  return new Promise(
    (resolve: (result: ListAdminSubscriptionsResult) => void, reject): void => {
      client.listAdminSubscriptions(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface GetAdminSubscriptionResult {
  request: GetAdminSubscriptionRequest;
  response: GetAdminSubscriptionResponse;
}

export function getSubscription(
  request: GetAdminSubscriptionRequest
): Promise<GetAdminSubscriptionResult> {
  return new Promise(
    (resolve: (result: GetAdminSubscriptionResult) => void, reject): void => {
      client.getAdminSubscription(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface UpdateAdminSubscriptionResult {
  request: UpdateAdminSubscriptionRequest;
  response: UpdateAdminSubscriptionResponse;
}

export function updateSubscription(
  request: UpdateAdminSubscriptionRequest
): Promise<UpdateAdminSubscriptionResult> {
  return new Promise(
    (
      resolve: (result: UpdateAdminSubscriptionResult) => void,
      reject
    ): void => {
      client.updateAdminSubscription(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface DeleteAdminSubscriptionResult {
  request: DeleteAdminSubscriptionRequest;
  response: DeleteAdminSubscriptionResponse;
}

export function deleteSubscription(
  request: DeleteAdminSubscriptionRequest
): Promise<DeleteAdminSubscriptionResult> {
  return new Promise(
    (
      resolve: (result: DeleteAdminSubscriptionResult) => void,
      reject
    ): void => {
      client.deleteAdminSubscription(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface EnableAdminSubscriptionResult {
  request: EnableAdminSubscriptionRequest;
  response: EnableAdminSubscriptionResponse;
}

export function enableSubscription(
  request: EnableAdminSubscriptionRequest
): Promise<EnableAdminSubscriptionResult> {
  return new Promise(
    (
      resolve: (result: EnableAdminSubscriptionResult) => void,
      reject
    ): void => {
      client.enableAdminSubscription(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface DisableAdminSubscriptionResult {
  request: DisableAdminSubscriptionRequest;
  response: DisableAdminSubscriptionResponse;
}

export function disableSubscription(
  request: DisableAdminSubscriptionRequest
): Promise<DisableAdminSubscriptionResult> {
  return new Promise(
    (
      resolve: (result: DisableAdminSubscriptionResult) => void,
      reject
    ): void => {
      client.disableAdminSubscription(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new NotificationServiceError(
                isNotNull(error) ? error.message : undefined,
                request,
                error
              )
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}
