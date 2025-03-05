import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  GetAuthenticationURLRequest,
  GetAuthenticationURLResponse,
  ExchangeTokenRequest,
  ExchangeTokenResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  SignInRequest,
  SignInResponse,
  SwitchOrganizationRequest,
  SwitchOrganizationResponse
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

export interface SignInResult {
  request: SignInRequest;
  response: SignInResponse;
}

export function signIn(request: SignInRequest): Promise<SignInResult> {
  return new Promise(
    (resolve: (result: SignInResult) => void, reject): void => {
      client.signIn(request, getMetaData(), (error, response): void => {
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

export interface SwitchOrganizationResult {
  request: SwitchOrganizationRequest;
  response: SwitchOrganizationResponse;
}

export function switchOrganization(
  request: SwitchOrganizationRequest
): Promise<SwitchOrganizationResult> {
  return new Promise(
    (resolve: (result: SwitchOrganizationResult) => void, reject): void => {
      client.switchOrganization(
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
