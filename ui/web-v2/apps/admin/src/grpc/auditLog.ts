import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ListAdminAuditLogsRequest,
  ListAdminAuditLogsResponse,
  ListAuditLogsRequest,
  ListAuditLogsResponse,
  ListFeatureHistoryRequest,
  ListFeatureHistoryResponse,
} from '../proto/auditlog/service_pb';
import {
  AuditLogServiceClient,
  ServiceError,
} from '../proto/auditlog/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class AuditLogServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuditLogServiceError);
    }
    this.name = 'AuditLogServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AuditLogServiceClient(urls.GRPC);

export interface ListAuditLogsResult {
  request: ListAuditLogsRequest;
  response: ListAuditLogsResponse;
}

export function listAuditLogs(
  request: ListAuditLogsRequest
): Promise<ListAuditLogsResult> {
  return new Promise(
    (resolve: (result: ListAuditLogsResult) => void, reject): void => {
      client.listAuditLogs(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new AuditLogServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListAdminAuditLogsResult {
  request: ListAdminAuditLogsRequest;
  response: ListAdminAuditLogsResponse;
}

export function listAdminAuditLogs(
  request: ListAdminAuditLogsRequest
): Promise<ListAdminAuditLogsResult> {
  return new Promise(
    (resolve: (result: ListAdminAuditLogsResult) => void, reject): void => {
      client.listAdminAuditLogs(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AuditLogServiceError(
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

export interface ListFeatureHistoryResult {
  request: ListFeatureHistoryRequest;
  response: ListFeatureHistoryResponse;
}

export function listFeatureHistory(
  request: ListFeatureHistoryRequest
): Promise<ListFeatureHistoryResult> {
  return new Promise(
    (resolve: (result: ListFeatureHistoryResult) => void, reject): void => {
      client.listFeatureHistory(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new AuditLogServiceError(
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
