import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateEnvironmentRequest,
  CreateEnvironmentResponse,
  DeleteEnvironmentRequest,
  DeleteEnvironmentResponse,
  GetEnvironmentRequest,
  GetEnvironmentResponse,
  ListEnvironmentsRequest,
  ListEnvironmentsResponse,
  UpdateEnvironmentRequest,
  UpdateEnvironmentResponse,
} from '../proto/environment/service_pb';
import {
  EnvironmentServiceClient,
  ServiceError,
} from '../proto/environment/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class EnvironmentServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, EnvironmentServiceError);
    }
    this.name = 'EnvironmentServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new EnvironmentServiceClient(urls.GRPC);

export interface GetEnvironmentResult {
  request: GetEnvironmentRequest;
  response: GetEnvironmentResponse;
}

export function getEnvironment(
  request: GetEnvironmentRequest
): Promise<GetEnvironmentResult> {
  return new Promise(
    (resolve: (result: GetEnvironmentResult) => void, reject): void => {
      client.getEnvironment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new EnvironmentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListEnvironmentsResult {
  request: ListEnvironmentsRequest;
  response: ListEnvironmentsResponse;
}

export function listEnvironments(
  request: ListEnvironmentsRequest
): Promise<ListEnvironmentsResult> {
  return new Promise(
    (resolve: (result: ListEnvironmentsResult) => void, reject): void => {
      client.listEnvironments(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new EnvironmentServiceError(
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

export interface CreateEnvironmentResult {
  request: CreateEnvironmentRequest;
  response: CreateEnvironmentResponse;
}

export function createEnvironment(
  request: CreateEnvironmentRequest
): Promise<CreateEnvironmentResult> {
  return new Promise(
    (resolve: (result: CreateEnvironmentResult) => void, reject): void => {
      client.createEnvironment(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new EnvironmentServiceError(
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

export interface DeleteEnvironmentResult {
  request: DeleteEnvironmentRequest;
  response: DeleteEnvironmentResponse;
}

export function deleteEnvironment(
  request: DeleteEnvironmentRequest
): Promise<DeleteEnvironmentResult> {
  return new Promise(
    (resolve: (result: DeleteEnvironmentResult) => void, reject): void => {
      client.deleteEnvironment(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new EnvironmentServiceError(
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

export interface UpdateEnvironmentResult {
  request: UpdateEnvironmentRequest;
  response: UpdateEnvironmentResponse;
}

export function updateEnvironment(
  request: UpdateEnvironmentRequest
): Promise<UpdateEnvironmentResult> {
  return new Promise(
    (resolve: (result: UpdateEnvironmentResult) => void, reject): void => {
      client.updateEnvironment(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new EnvironmentServiceError(
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
