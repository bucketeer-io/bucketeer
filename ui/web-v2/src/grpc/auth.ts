import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  GetAuthenticationURLRequest,
  GetAuthenticationURLResponse,
  ExchangeBucketeerTokenRequest,
  ExchangeBucketeerTokenResponse,
  RefreshBucketeerTokenRequest,
  RefreshBucketeerTokenResponse
} from '../proto/auth/service_pb';
import {
  AuthServiceClient,
  ServiceError
} from '../proto/auth/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class AuthServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuthServiceError);
    }
    this.name = 'AuthServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AuthServiceClient(urls.GRPC);

export interface GetAuthenticationResult {
  request: GetAuthenticationURLRequest;
  response: GetAuthenticationURLResponse;
}

export function getAuthenticationURL(
  request: GetAuthenticationURLRequest
): Promise<GetAuthenticationResult> {
  return new Promise(
    (resolve: (result: GetAuthenticationResult) => void, reject): void => {
      client.getAuthenticationURL(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AuthServiceError(extractErrorMessage(error), request, error)
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface ExchangeBucketeerTokenResult {
  request: ExchangeBucketeerTokenRequest;
  response: ExchangeBucketeerTokenResponse;
}

export function exchangeBucketeerToken(
  request: ExchangeBucketeerTokenRequest
): Promise<ExchangeBucketeerTokenResult> {
  return new Promise(
    (resolve: (result: ExchangeBucketeerTokenResult) => void, reject): void => {
      client.exchangeBucketeerToken(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AuthServiceError(extractErrorMessage(error), request, error)
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}

export interface RefreshBucketeerTokenResult {
  request: RefreshBucketeerTokenRequest;
  response: RefreshBucketeerTokenResponse;
}

export function refreshBucketeerToken(
  request: RefreshBucketeerTokenRequest
): Promise<RefreshBucketeerTokenResult> {
  return new Promise(
    (resolve: (result: RefreshBucketeerTokenResult) => void, reject): void => {
      client.refreshBucketeerToken(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AuthServiceError(extractErrorMessage(error), request, error)
            );
          } else {
            resolve({ request, response });
          }
        }
      );
    }
  );
}
