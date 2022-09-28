import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';
import { Undefinable } from 'option-t/lib/Undefinable/Undefinable';

import { urls } from '../config';
import {
  CreateSubscriptionRequest,
  CreateSubscriptionResponse,
  DeleteSubscriptionRequest,
  DeleteSubscriptionResponse,
  DisableSubscriptionRequest,
  DisableSubscriptionResponse,
  EnableSubscriptionRequest,
  EnableSubscriptionResponse,
  ListSubscriptionsRequest,
  ListSubscriptionsResponse,
  UpdateSubscriptionRequest,
  UpdateSubscriptionResponse,
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

export interface CreateSubscriptionResult {
  request: CreateSubscriptionRequest;
  response: CreateSubscriptionResponse;
}

export function createSubscription(
  request: CreateSubscriptionRequest
): Promise<CreateSubscriptionResult> {
  return new Promise(
    (resolve: (result: CreateSubscriptionResult) => void, reject): void => {
      client.createSubscription(
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

export interface ListSubscriptionsResult {
  request: ListSubscriptionsRequest;
  response: ListSubscriptionsResponse;
}

export function listSubscriptions(
  request: ListSubscriptionsRequest
): Promise<ListSubscriptionsResult> {
  return new Promise(
    (resolve: (result: ListSubscriptionsResult) => void, reject): void => {
      client.listSubscriptions(
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

export interface UpdateSubscriptionResult {
  request: UpdateSubscriptionRequest;
  response: UpdateSubscriptionResponse;
}

export function updateSubscription(
  request: UpdateSubscriptionRequest
): Promise<UpdateSubscriptionResult> {
  return new Promise(
    (resolve: (result: UpdateSubscriptionResult) => void, reject): void => {
      client.updateSubscription(
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

export interface DeleteSubscriptionResult {
  request: DeleteSubscriptionRequest;
  response: DeleteSubscriptionResponse;
}

export function deleteSubscription(
  request: DeleteSubscriptionRequest
): Promise<DeleteSubscriptionResult> {
  return new Promise(
    (resolve: (result: DeleteSubscriptionResult) => void, reject): void => {
      client.deleteSubscription(
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

export interface EnableSubscriptionResult {
  request: EnableSubscriptionRequest;
  response: EnableSubscriptionResponse;
}

export function enableSubscription(
  request: EnableSubscriptionRequest
): Promise<EnableSubscriptionResult> {
  return new Promise(
    (resolve: (result: EnableSubscriptionResult) => void, reject): void => {
      client.enableSubscription(
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

export interface DisableSubscriptionResult {
  request: DisableSubscriptionRequest;
  response: DisableSubscriptionResponse;
}

export function disableSubscription(
  request: DisableSubscriptionRequest
): Promise<DisableSubscriptionResult> {
  return new Promise(
    (resolve: (result: DisableSubscriptionResult) => void, reject): void => {
      client.disableSubscription(
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
