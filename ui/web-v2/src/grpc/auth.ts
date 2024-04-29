import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ExchangeTokenRequest,
  ExchangeTokenResponse,
  GetAuthCodeURLRequest,
  GetAuthCodeURLResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
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

export interface GetAuthCodeURLResult {
  request: GetAuthCodeURLRequest;
  response: GetAuthCodeURLResponse;
}

export function getAuthCodeURL(
  request: GetAuthCodeURLRequest
): Promise<GetAuthCodeURLResult> {
  return new Promise(
    (resolve: (result: GetAuthCodeURLResult) => void, reject): void => {
      client.getAuthCodeURL(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AuthServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

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

export interface ExchangeTokenResult {
  request: ExchangeTokenRequest;
  response: ExchangeTokenResponse;
}

export function exchangeToken(
  request: ExchangeTokenRequest
): Promise<ExchangeTokenResult> {
  return new Promise(
    (resolve: (result: ExchangeTokenResult) => void, reject): void => {
      client.exchangeToken(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AuthServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
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

export interface RefreshTokenResult {
  request: RefreshTokenRequest;
  response: RefreshTokenResponse;
}

export function refreshToken(
  request: RefreshTokenRequest
): Promise<RefreshTokenResult> {
  return new Promise(
    (resolve: (result: RefreshTokenResult) => void, reject): void => {
      client.refreshToken(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AuthServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
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
