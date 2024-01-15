import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ChangeAccountRoleRequest,
  ChangeAccountRoleResponse,
  ConvertAccountRequest,
  ConvertAccountResponse,
  CreateAccountRequest,
  CreateAccountResponse,
  DisableAccountRequest,
  DisableAccountResponse,
  EnableAccountRequest,
  EnableAccountResponse,
  GetAccountRequest,
  GetAccountResponse,
  GetMeRequest,
  GetMeResponse,
  GetMeV2Request,
  GetMeV2Response,
  GetMyOrganizationsRequest,
  GetMyOrganizationsResponse,
  ListAccountsRequest,
  ListAccountsResponse,
} from '../proto/account/service_pb';
import {
  AccountServiceClient,
  ServiceError,
} from '../proto/account/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

const client = new AccountServiceClient(urls.GRPC);

export interface GetMeResult {
  request: GetMeRequest;
  response: GetMeResponse;
}

export function getMe(request: GetMeRequest): Promise<GetMeResult> {
  return new Promise((resolve: (result: GetMeResult) => void, reject): void => {
    client.getMe(request, getMetaData(), (error, response): void => {
      if (isNotNull(error) || isNull(response)) {
        reject(
          new AccountServiceError(extractErrorMessage(error), request, error)
        );
      } else {
        resolve({ request, response });
      }
    });
  });
}

export interface GetMyOrganizationsResult {
  request: GetMyOrganizationsRequest;
  response: GetMyOrganizationsResponse;
}

export function getMyOrganizations(request: GetMyOrganizationsRequest): Promise<GetMyOrganizationsResult> {
  return new Promise((resolve: (result: GetMyOrganizationsResult) => void, reject): void => {
    client.getMyOrganizations(request, getMetaData(), (error, response): void => {
      if (isNotNull(error) || isNull(response)) {
        reject(
          new AccountServiceError(extractErrorMessage(error), request, error)
        );
      } else {
        resolve({ request, response });
      }
    });
  });
}

export interface GetMeV2Result {
  request: GetMeV2Request;
  response: GetMeV2Response;
}

export function getMeV2(request: GetMeV2Request): Promise<GetMeV2Result> {
  return new Promise((resolve: (result: GetMeV2Result) => void, reject): void => {
    client.getMeV2(request, getMetaData(), (error, response): void => {
      if (isNotNull(error) || isNull(response)) {
        reject(
          new AccountServiceError(extractErrorMessage(error), request, error)
        );
      } else {
        resolve({ request, response });
      }
    });
  });
}

export interface CreateAccountResult {
  request: CreateAccountRequest;
  response: CreateAccountResponse;
}

export function createAccount(
  request: CreateAccountRequest
): Promise<CreateAccountResult> {
  return new Promise(
    (resolve: (result: CreateAccountResult) => void, reject): void => {
      client.createAccount(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AccountServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface EnableAccountResult {
  request: EnableAccountRequest;
  response: EnableAccountResponse;
}

export function enableAccount(
  request: EnableAccountRequest
): Promise<EnableAccountResult> {
  return new Promise(
    (resolve: (result: EnableAccountResult) => void, reject): void => {
      client.enableAccount(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AccountServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DisableAccountResult {
  request: DisableAccountRequest;
  response: DisableAccountResponse;
}

export function disableAccount(
  request: DisableAccountRequest
): Promise<DisableAccountResult> {
  return new Promise(
    (resolve: (result: DisableAccountResult) => void, reject): void => {
      client.disableAccount(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AccountServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ChangeAccountRoleResult {
  request: ChangeAccountRoleRequest;
  response: ChangeAccountRoleResponse;
}

export function changeAccountRole(
  request: ChangeAccountRoleRequest
): Promise<ChangeAccountRoleResult> {
  return new Promise(
    (resolve: (result: ChangeAccountRoleResult) => void, reject): void => {
      client.changeAccountRole(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AccountServiceError(
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

export interface ConvertAccountResult {
  request: ConvertAccountRequest;
  response: ConvertAccountResponse;
}

export function convertAccount(
  request: ConvertAccountRequest
): Promise<ConvertAccountResult> {
  return new Promise(
    (resolve: (result: ConvertAccountResult) => void, reject): void => {
      client.convertAccount(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AccountServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface GetAccountResult {
  request: GetAccountRequest;
  response: GetAccountResponse;
}

export function getAccount(
  request: GetAccountRequest
): Promise<GetAccountResult> {
  return new Promise(
    (resolve: (result: GetAccountResult) => void, reject): void => {
      client.getAccount(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AccountServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListAccountsResult {
  request: ListAccountsRequest;
  response: ListAccountsResponse;
}

export function listAccounts(
  request: ListAccountsRequest
): Promise<ListAccountsResult> {
  return new Promise(
    (resolve: (result: ListAccountsResult) => void, reject): void => {
      client.listAccounts(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AccountServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export class AccountServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AccountServiceError);
    }
    this.name = 'AccountServiceError';
    this.request = request;
    this.error = error;
  }
}
