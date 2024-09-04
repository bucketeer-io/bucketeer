import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateAccountV2Request,
  CreateAccountV2Response,
  DisableAccountV2Request,
  DisableAccountV2Response,
  EnableAccountV2Request,
  EnableAccountV2Response,
  GetAccountV2Request,
  GetAccountV2Response,
  GetMeRequest,
  GetMeResponse,
  GetMyOrganizationsRequest,
  GetMyOrganizationsResponse,
  ListAccountsV2Request,
  ListAccountsV2Response,
  UpdateAccountV2Request,
  UpdateAccountV2Response,
  CreateSearchFilterRequest,
  CreateSearchFilterResponse,
  UpdateSearchFilterRequest,
  UpdateSearchFilterResponse,
  DeleteSearchFilterRequest,
  DeleteSearchFilterResponse
} from '../proto/account/service_pb';
import {
  AccountServiceClient,
  ServiceError
} from '../proto/account/service_pb_service';

import { extractErrorMessage } from './messages';
import {
  checkUnauthenticatedError,
  getMetaDataForClient as getMetaData
} from './utils';
import { UNAUTHENTICATED_ERROR } from '../middlewares/thunkErrorHandler';

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

export function getMyOrganizations(
  request: GetMyOrganizationsRequest
): Promise<GetMyOrganizationsResult> {
  return new Promise(
    (resolve: (result: GetMyOrganizationsResult) => void, reject): void => {
      client.getMyOrganizations(
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

export interface CreateAccountResult {
  request: CreateAccountV2Request;
  response: CreateAccountV2Response;
}

export function createAccount(
  request: CreateAccountV2Request
): Promise<CreateAccountResult> {
  return new Promise(
    (resolve: (result: CreateAccountResult) => void, reject): void => {
      client.createAccountV2(
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

export interface EnableAccountResult {
  request: EnableAccountV2Request;
  response: EnableAccountV2Response;
}

export function enableAccount(
  request: EnableAccountV2Request
): Promise<EnableAccountResult> {
  return new Promise(
    (resolve: (result: EnableAccountResult) => void, reject): void => {
      client.enableAccountV2(
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

export interface DisableAccountResult {
  request: DisableAccountV2Request;
  response: DisableAccountV2Response;
}

export function disableAccount(
  request: DisableAccountV2Request
): Promise<DisableAccountResult> {
  return new Promise(
    (resolve: (result: DisableAccountResult) => void, reject): void => {
      client.disableAccountV2(
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

export interface UpdateAccountResult {
  request: UpdateAccountV2Request;
  response: UpdateAccountV2Response;
}

export function updateAccount(
  request: UpdateAccountV2Request
): Promise<UpdateAccountResult> {
  return new Promise(
    (resolve: (result: UpdateAccountResult) => void, reject): void => {
      client.updateAccountV2(
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

export interface GetAccountResult {
  request: GetAccountV2Request;
  response: GetAccountV2Response;
}

export function getAccount(
  request: GetAccountV2Request
): Promise<GetAccountResult> {
  return new Promise(
    (resolve: (result: GetAccountResult) => void, reject): void => {
      client.getAccountV2(request, getMetaData(), (error, response): void => {
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
  request: ListAccountsV2Request;
  response: ListAccountsV2Response;
}

export function listAccounts(
  request: ListAccountsV2Request
): Promise<ListAccountsResult> {
  return new Promise(
    (resolve: (result: ListAccountsResult) => void, reject): void => {
      client.listAccountsV2(request, getMetaData(), (error, response): void => {
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

export interface createSearchFilterResult {
  request: CreateSearchFilterRequest;
  response: CreateSearchFilterResponse;
}

export function createSearchFilter(
  request: CreateSearchFilterRequest
): Promise<createSearchFilterResult> {
  return new Promise(
    (resolve: (result: createSearchFilterResult) => void, reject): void => {
      client.createSearchFilter(
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

export interface updateSearchFilterResult {
  request: UpdateSearchFilterRequest;
  response: UpdateSearchFilterResponse;
}

export function updateSearchFilter(
  request: UpdateSearchFilterRequest
): Promise<updateSearchFilterResult> {
  return new Promise(
    (resolve: (result: updateSearchFilterResult) => void, reject): void => {
      client.updateSearchFilter(
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

export interface deleteSearchFilterResult {
  request: DeleteSearchFilterRequest;
  response: DeleteSearchFilterResponse;
}

export function deleteSearchFilter(
  request: DeleteSearchFilterRequest
): Promise<deleteSearchFilterResult> {
  return new Promise(
    (resolve: (result: deleteSearchFilterResult) => void, reject): void => {
      client.deleteSearchFilter(
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

export class AccountServiceError<Request> extends Error {
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
      Error.captureStackTrace(this, AccountServiceError);
    }
    this.name = 'AccountServiceError';
    this.request = request;
    this.error = error;
  }
}
