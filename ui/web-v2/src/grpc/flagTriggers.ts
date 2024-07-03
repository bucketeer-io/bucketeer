import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateFlagTriggerRequest,
  CreateFlagTriggerResponse,
  ListFlagTriggersRequest,
  ListFlagTriggersResponse,
  DeleteFlagTriggerRequest,
  DeleteFlagTriggerResponse,
  UpdateFlagTriggerRequest,
  UpdateFlagTriggerResponse,
  ResetFlagTriggerRequest,
  ResetFlagTriggerResponse,
  EnableFlagTriggerRequest,
  EnableFlagTriggerResponse,
  DisableFlagTriggerRequest,
  DisableFlagTriggerResponse
} from '../proto/feature/service_pb';
import {
  FeatureServiceClient,
  ServiceError
} from '../proto/feature/service_pb_service';

import { extractErrorMessage } from './messages';
import {
  checkUnauthenticatedError,
  getMetaDataForClient as getMetaData
} from './utils';
import { UNAUTHENTICATED_ERROR } from '../middlewares/thunkErrorHandler';

export class FlagTriggersServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    if (checkUnauthenticatedError(error.code)) {
      super(UNAUTHENTICATED_ERROR);
    } else {
      super(message);
    }
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, FlagTriggersServiceError);
    }
    this.name = 'FlagTriggersServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new FeatureServiceClient(urls.GRPC);

export interface CreateFlagTriggerResult {
  request: CreateFlagTriggerRequest;
  response: CreateFlagTriggerResponse;
}

export function createFlagTrigger(
  request: CreateFlagTriggerRequest
): Promise<CreateFlagTriggerResult> {
  return new Promise(
    (resolve: (result: CreateFlagTriggerResult) => void, reject): void => {
      client.createFlagTrigger(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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

export interface ListFlagTriggerResult {
  request: ListFlagTriggersRequest;
  response: ListFlagTriggersResponse;
}

export function listFlagTriggers(
  request: ListFlagTriggersRequest
): Promise<ListFlagTriggerResult> {
  return new Promise(
    (resolve: (result: ListFlagTriggerResult) => void, reject): void => {
      client.listFlagTriggers(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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

export interface UpdateFlagTriggerResult {
  request: UpdateFlagTriggerRequest;
  response: UpdateFlagTriggerResponse;
}

export function updateFlagTrigger(
  request: UpdateFlagTriggerRequest
): Promise<UpdateFlagTriggerResult> {
  return new Promise(
    (resolve: (result: UpdateFlagTriggerResult) => void, reject): void => {
      client.updateFlagTrigger(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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

export interface DeleteFlagTriggerResult {
  request: DeleteFlagTriggerRequest;
  response: DeleteFlagTriggerResponse;
}

export function deleteFlagTrigger(
  request: DeleteFlagTriggerRequest
): Promise<DeleteFlagTriggerResult> {
  return new Promise(
    (resolve: (result: DeleteFlagTriggerResult) => void, reject): void => {
      client.deleteFlagTrigger(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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

export interface ResetFlagTriggerResult {
  request: ResetFlagTriggerRequest;
  response: ResetFlagTriggerResponse;
}

export function resetFlagTrigger(
  request: ResetFlagTriggerRequest
): Promise<ResetFlagTriggerResult> {
  return new Promise(
    (resolve: (result: ResetFlagTriggerResult) => void, reject): void => {
      client.resetFlagTrigger(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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

export interface EnableFlagTriggerResult {
  request: EnableFlagTriggerRequest;
  response: EnableFlagTriggerResponse;
}

export function enableFlagTrigger(
  request: EnableFlagTriggerRequest
): Promise<EnableFlagTriggerResult> {
  return new Promise(
    (resolve: (result: EnableFlagTriggerResult) => void, reject): void => {
      client.enableFlagTrigger(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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

export interface DisableFlagTriggerResult {
  request: DisableFlagTriggerRequest;
  response: DisableFlagTriggerResponse;
}

export function disableFlagTrigger(
  request: DisableFlagTriggerRequest
): Promise<DisableFlagTriggerResult> {
  return new Promise(
    (resolve: (result: DisableFlagTriggerResult) => void, reject): void => {
      client.disableFlagTrigger(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FlagTriggersServiceError(
                extractErrorMessage(error),
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
