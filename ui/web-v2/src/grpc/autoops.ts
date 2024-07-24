import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateAutoOpsRuleRequest,
  CreateAutoOpsRuleResponse,
  DeleteAutoOpsRuleRequest,
  DeleteAutoOpsRuleResponse,
  ListAutoOpsRulesRequest,
  ListAutoOpsRulesResponse,
  ListOpsCountsRequest,
  ListOpsCountsResponse,
  UpdateAutoOpsRuleRequest,
  UpdateAutoOpsRuleResponse,
  StopAutoOpsRuleRequest,
  StopAutoOpsRuleResponse
} from '../proto/autoops/service_pb';
import {
  AutoOpsServiceClient,
  ServiceError
} from '../proto/autoops/service_pb_service';

import { extractErrorMessage } from './messages';
import {
  checkUnauthenticatedError,
  getMetaDataForClient as getMetaData
} from './utils';
import { UNAUTHENTICATED_ERROR } from '../middlewares/thunkErrorHandler';

export class AutoOpsServiceError<Request> extends Error {
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
      Error.captureStackTrace(this, AutoOpsServiceError);
    }
    this.name = 'AutoOpsServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AutoOpsServiceClient(urls.GRPC);

export interface CreateAutoOpsRuleResult {
  request: CreateAutoOpsRuleRequest;
  response: CreateAutoOpsRuleResponse;
}

export function createAutoOpsRule(
  request: CreateAutoOpsRuleRequest
): Promise<CreateAutoOpsRuleResult> {
  return new Promise(
    (resolve: (result: CreateAutoOpsRuleResult) => void, reject): void => {
      client.createAutoOpsRule(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AutoOpsServiceError(
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

export interface ListAutoOpsRulesResult {
  request: ListAutoOpsRulesRequest;
  response: ListAutoOpsRulesResponse;
}

export function listAutoOpsRules(
  request: ListAutoOpsRulesRequest
): Promise<ListAutoOpsRulesResult> {
  return new Promise(
    (resolve: (result: ListAutoOpsRulesResult) => void, reject): void => {
      client.listAutoOpsRules(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AutoOpsServiceError(
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

export interface UpdateAutoOpsRuleResult {
  request: UpdateAutoOpsRuleRequest;
  response: UpdateAutoOpsRuleResponse;
}

export function updateAutoOpsRule(
  request: UpdateAutoOpsRuleRequest
): Promise<UpdateAutoOpsRuleResult> {
  return new Promise(
    (resolve: (result: UpdateAutoOpsRuleResult) => void, reject): void => {
      client.updateAutoOpsRule(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AutoOpsServiceError(
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

export interface DeleteAutoOpsRuleResult {
  request: DeleteAutoOpsRuleRequest;
  response: DeleteAutoOpsRuleResponse;
}

export function deleteAutoOpsRule(
  request: DeleteAutoOpsRuleRequest
): Promise<DeleteAutoOpsRuleResult> {
  return new Promise(
    (resolve: (result: DeleteAutoOpsRuleResult) => void, reject): void => {
      client.deleteAutoOpsRule(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AutoOpsServiceError(
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

export interface StopAutoOpsRuleResult {
  request: StopAutoOpsRuleRequest;
  response: StopAutoOpsRuleResponse;
}

export function stopAutoOpsRule(
  request: StopAutoOpsRuleRequest
): Promise<StopAutoOpsRuleResult> {
  return new Promise(
    (resolve: (result: StopAutoOpsRuleResult) => void, reject): void => {
      client.stopAutoOpsRule(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AutoOpsServiceError(
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

export interface ListOpsCountsResult {
  request: ListOpsCountsRequest;
  response: ListOpsCountsResponse;
}

export function listOpsCounts(
  request: ListOpsCountsRequest
): Promise<ListOpsCountsResult> {
  return new Promise(
    (resolve: (result: ListOpsCountsResult) => void, reject): void => {
      client.listOpsCounts(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AutoOpsServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}
