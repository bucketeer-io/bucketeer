import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateWebhookRequest,
  CreateWebhookResponse,
  DeleteWebhookRequest,
  DeleteWebhookResponse,
  ListWebhooksRequest,
  ListWebhooksResponse,
  UpdateWebhookRequest,
  UpdateWebhookResponse,
  GetWebhookRequest,
  GetWebhookResponse,
} from '../proto/autoops/service_pb';
import {
  ServiceError,
  AutoOpsServiceClient,
} from '../proto/autoops/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class AutoOpsServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AutoOpsServiceError);
    }
    this.name = 'AutoOpsServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AutoOpsServiceClient(urls.GRPC);

export interface CreateWebhookResult {
  request: CreateWebhookRequest;
  response: CreateWebhookResponse;
}

export function createWebhook(
  request: CreateWebhookRequest
): Promise<CreateWebhookResult> {
  return new Promise(
    (resolve: (result: CreateWebhookResult) => void, reject): void => {
      client.createWebhook(request, getMetaData(), (error, response): void => {
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

export interface ListWebhooksResult {
  request: ListWebhooksRequest;
  response: ListWebhooksResponse;
}

export function listWebhooks(
  request: ListWebhooksRequest
): Promise<ListWebhooksResult> {
  return new Promise(
    (resolve: (result: ListWebhooksResult) => void, reject): void => {
      client.listWebhooks(request, getMetaData(), (error, response): void => {
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

export interface GetWebhooksResult {
  request: GetWebhookRequest;
  response: GetWebhookResponse;
}

export function getWebhook(
  request: GetWebhookRequest
): Promise<GetWebhooksResult> {
  return new Promise(
    (resolve: (result: GetWebhooksResult) => void, reject): void => {
      client.getWebhook(request, getMetaData(), (error, response): void => {
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

export interface UpdateWebhookResult {
  request: UpdateWebhookRequest;
  response: UpdateWebhookResponse;
}

export function updateWebhook(
  request: UpdateWebhookRequest
): Promise<UpdateWebhookResult> {
  return new Promise(
    (resolve: (result: UpdateWebhookResult) => void, reject): void => {
      client.updateWebhook(request, getMetaData(), (error, response): void => {
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

export interface DeleteWebhookResult {
  request: DeleteWebhookRequest;
  response: DeleteWebhookResponse;
}

export function deleteWebhook(
  request: DeleteWebhookRequest
): Promise<DeleteWebhookResult> {
  return new Promise(
    (resolve: (result: DeleteWebhookResult) => void, reject): void => {
      client.deleteWebhook(request, getMetaData(), (error, response): void => {
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
