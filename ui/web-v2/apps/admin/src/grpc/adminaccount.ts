import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateAdminAccountRequest,
  CreateAdminAccountResponse,
  DisableAdminAccountRequest,
  DisableAdminAccountResponse,
  EnableAdminAccountRequest,
  EnableAdminAccountResponse,
  GetAdminAccountRequest,
  GetAdminAccountResponse,
  ListAdminAccountsRequest,
  ListAdminAccountsResponse,
} from '../proto/account/service_pb';
import {
  AccountServiceClient,
  ServiceError,
} from '../proto/account/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class AdminAccountServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AdminAccountServiceError);
    }
    this.name = 'AdminAccountServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AccountServiceClient(urls.GRPC);

export interface CreateAdminAccountResult {
  request: CreateAdminAccountRequest;
  response: CreateAdminAccountResponse;
}

export function createAdminAccount(
  request: CreateAdminAccountRequest
): Promise<CreateAdminAccountResult> {
  return new Promise(
    (resolve: (result: CreateAdminAccountResult) => void, reject): void => {
      client.createAdminAccount(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AdminAccountServiceError(
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

export interface EnableAdminAccountResult {
  request: EnableAdminAccountRequest;
  response: EnableAdminAccountResponse;
}

export function enableAdminAccount(
  request: EnableAdminAccountRequest
): Promise<EnableAdminAccountResult> {
  return new Promise(
    (resolve: (result: EnableAdminAccountResult) => void, reject): void => {
      client.enableAdminAccount(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AdminAccountServiceError(
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

export interface DisableAdminAccountResult {
  request: DisableAdminAccountRequest;
  response: DisableAdminAccountResponse;
}

export function disableAdminAccount(
  request: DisableAdminAccountRequest
): Promise<DisableAdminAccountResult> {
  return new Promise(
    (resolve: (result: DisableAdminAccountResult) => void, reject): void => {
      client.disableAdminAccount(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AdminAccountServiceError(
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

export interface GetAdminAccountResult {
  request: GetAdminAccountRequest;
  response: GetAdminAccountResponse;
}

export function getAdminAccount(
  request: GetAdminAccountRequest
): Promise<GetAdminAccountResult> {
  return new Promise(
    (resolve: (result: GetAdminAccountResult) => void, reject): void => {
      client.getAdminAccount(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AdminAccountServiceError(
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

export interface ListAdminAccountsResult {
  request: ListAdminAccountsRequest;
  response: ListAdminAccountsResponse;
}

export function listAdminAccounts(
  request: ListAdminAccountsRequest
): Promise<ListAdminAccountsResult> {
  return new Promise(
    (resolve: (result: ListAdminAccountsResult) => void, reject): void => {
      client.listAdminAccounts(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AdminAccountServiceError(
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
