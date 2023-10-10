import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ArchiveEnvironmentV2Request,
  ArchiveEnvironmentV2Response,
  CreateEnvironmentV2Request,
  CreateEnvironmentV2Response,
  GetEnvironmentV2Request,
  GetEnvironmentV2Response,
  ListEnvironmentsV2Request,
  ListEnvironmentsV2Response,
  UnarchiveEnvironmentV2Request,
  UnarchiveEnvironmentV2Response,
  UpdateEnvironmentV2Request,
  UpdateEnvironmentV2Response,
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
  request: GetEnvironmentV2Request;
  response: GetEnvironmentV2Response;
}

export function getEnvironment(
  request: GetEnvironmentV2Request
): Promise<GetEnvironmentResult> {
  return new Promise(
    (resolve: (result: GetEnvironmentResult) => void, reject): void => {
      client.getEnvironmentV2(request, getMetaData(), (error, response): void => {
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
  request: ListEnvironmentsV2Request;
  response: ListEnvironmentsV2Response;
}

export function listEnvironments(
  request: ListEnvironmentsV2Request
): Promise<ListEnvironmentsResult> {
  return new Promise(
    (resolve: (result: ListEnvironmentsResult) => void, reject): void => {
      client.listEnvironmentsV2(
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
  request: CreateEnvironmentV2Request;
  response: CreateEnvironmentV2Response;
}

export function createEnvironment(
  request: CreateEnvironmentV2Request
): Promise<CreateEnvironmentResult> {
  return new Promise(
    (resolve: (result: CreateEnvironmentResult) => void, reject): void => {
      client.createEnvironmentV2(
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

export interface ArchiveEnvironmentResult {
  request: ArchiveEnvironmentV2Request;
  response: ArchiveEnvironmentV2Response;
}

export function archiveEnvironment(
  request: ArchiveEnvironmentV2Request
): Promise<ArchiveEnvironmentResult> {
  return new Promise(
    (resolve: (result: ArchiveEnvironmentResult) => void, reject): void => {
      client.archiveEnvironmentV2(
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

export interface UnarchiveEnvironmentResult {
  request: UnarchiveEnvironmentV2Request;
  response: UnarchiveEnvironmentV2Response;
}

export function unarchiveEnvironment(
  request: UnarchiveEnvironmentV2Request
): Promise<UnarchiveEnvironmentResult> {
  return new Promise(
    (resolve: (result: UnarchiveEnvironmentResult) => void, reject): void => {
      client.unarchiveEnvironmentV2(
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
  request: UpdateEnvironmentV2Request;
  response: UpdateEnvironmentV2Response;
}

export function updateEnvironment(
  request: UpdateEnvironmentV2Request
): Promise<UpdateEnvironmentResult> {
  return new Promise(
    (resolve: (result: UpdateEnvironmentResult) => void, reject): void => {
      client.updateEnvironmentV2(
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
