import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ChangeAPIKeyNameRequest,
  ChangeAPIKeyNameResponse,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  DisableAPIKeyRequest,
  DisableAPIKeyResponse,
  EnableAPIKeyRequest,
  EnableAPIKeyResponse,
  GetAPIKeyRequest,
  GetAPIKeyResponse,
  ListAPIKeysRequest,
  ListAPIKeysResponse,
} from '../proto/account/service_pb';
import {
  AccountServiceClient,
  ServiceError,
} from '../proto/account/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class APIKeyServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, APIKeyServiceError);
    }
    this.name = 'APIKeyServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AccountServiceClient(urls.GRPC);

export interface CreateAPIKeyResult {
  request: CreateAPIKeyRequest;
  response: CreateAPIKeyResponse;
}

export function createAPIKey(
  request: CreateAPIKeyRequest
): Promise<CreateAPIKeyResult> {
  return new Promise(
    (resolve: (result: CreateAPIKeyResult) => void, reject): void => {
      client.createAPIKey(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new APIKeyServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ChangeAPIKeyNameResult {
  request: ChangeAPIKeyNameRequest;
  response: ChangeAPIKeyNameResponse;
}

export function changeAPIKeyName(
  request: ChangeAPIKeyNameRequest
): Promise<ChangeAPIKeyNameResult> {
  return new Promise(
    (resolve: (result: ChangeAPIKeyNameResult) => void, reject): void => {
      client.changeAPIKeyName(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new APIKeyServiceError(extractErrorMessage(error), request, error)
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface GetAPIKeyResult {
  request: GetAPIKeyRequest;
  response: GetAPIKeyResponse;
}

export function getAPIKey(request: GetAPIKeyRequest): Promise<GetAPIKeyResult> {
  return new Promise(
    (resolve: (result: GetAPIKeyResult) => void, reject): void => {
      client.getAPIKey(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new APIKeyServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListAPIKeysResult {
  request: ListAPIKeysRequest;
  response: ListAPIKeysResponse;
}

export function listAPIKeys(
  request: ListAPIKeysRequest
): Promise<ListAPIKeysResult> {
  return new Promise(
    (resolve: (result: ListAPIKeysResult) => void, reject): void => {
      client.listAPIKeys(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new APIKeyServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface EnableAPIKeyResult {
  request: EnableAPIKeyRequest;
  response: EnableAPIKeyResponse;
}

export function enableAPIKey(
  request: EnableAPIKeyRequest
): Promise<EnableAPIKeyResult> {
  return new Promise(
    (resolve: (result: EnableAPIKeyResult) => void, reject): void => {
      client.enableAPIKey(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new APIKeyServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DisableAPIKeyResult {
  request: DisableAPIKeyRequest;
  response: DisableAPIKeyResponse;
}

export function disableAPIKey(
  request: DisableAPIKeyRequest
): Promise<DisableAPIKeyResult> {
  return new Promise(
    (resolve: (result: DisableAPIKeyResult) => void, reject): void => {
      client.disableAPIKey(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new APIKeyServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}
