import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreatePushRequest,
  CreatePushResponse,
  DeletePushRequest,
  DeletePushResponse,
  ListPushesRequest,
  ListPushesResponse,
  UpdatePushRequest,
  UpdatePushResponse,
} from '../proto/push/service_pb';
import {
  PushServiceClient,
  ServiceError,
} from '../proto/push/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class PushServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, PushServiceError);
    }
    this.name = 'PushServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new PushServiceClient(urls.GRPC);

export interface CreatePushResult {
  request: CreatePushRequest;
  response: CreatePushResponse;
}

export function createPush(
  request: CreatePushRequest
): Promise<CreatePushResult> {
  return new Promise(
    (resolve: (result: CreatePushResult) => void, reject): void => {
      client.createPush(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new PushServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListPushesResult {
  request: ListPushesRequest;
  response: ListPushesResponse;
}

export function listPushes(
  request: ListPushesRequest
): Promise<ListPushesResult> {
  return new Promise(
    (resolve: (result: ListPushesResult) => void, reject): void => {
      client.listPushes(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new PushServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UpdatePushResult {
  request: UpdatePushRequest;
  response: UpdatePushResponse;
}

export function updatePush(
  request: UpdatePushRequest
): Promise<UpdatePushResult> {
  return new Promise(
    (resolve: (result: UpdatePushResult) => void, reject): void => {
      client.updatePush(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new PushServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DeletePushResult {
  request: DeletePushRequest;
  response: DeletePushResponse;
}

export function deletePush(
  request: DeletePushRequest
): Promise<DeletePushResult> {
  return new Promise(
    (resolve: (result: DeletePushResult) => void, reject): void => {
      client.deletePush(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new PushServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}
