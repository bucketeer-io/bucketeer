import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ChangeAPIKeyNameRequest,
  ChangeAPIKeyNameResponse,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  DisableAPIKeyRequest,
  DisableAPIKeyResponse,
  EnableAPIKeyRequest,
  EnableAPIKeyResponse,
  GetAPIKeyRequest,
  GetAPIKeyResponse,
  ListAPIKeysRequest,
  ListAPIKeysResponse
} from '../proto/account/service_pb';
import {
  AccountServiceClient,
  ServiceError
} from '../proto/account/service_pb_service';
import { CodeReferenceServiceClient } from '../proto/coderef/service_pb_service';

import { extractErrorMessage } from './messages';
import {
  checkUnauthenticatedError,
  getMetaDataForClient as getMetaData
} from './utils';
import { UNAUTHENTICATED_ERROR } from '../middlewares/thunkErrorHandler';
import {
  GetCodeReferenceRequest,
  GetCodeReferenceResponse,
  ListCodeReferencesRequest,
  ListCodeReferencesResponse
} from '../proto/coderef/service_pb';

export class CodeRefsServiceError<Request> extends Error {
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
      Error.captureStackTrace(this, CodeRefsServiceError);
    }
    this.name = 'CodeRefsServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new CodeReferenceServiceClient(urls.GRPC);

export interface GetCodeReferenceResult {
  request: GetCodeReferenceRequest;
  response: GetCodeReferenceResponse;
}

export function getCodeRef(
  request: GetCodeReferenceRequest
): Promise<GetCodeReferenceResult> {
  return new Promise(
    (resolve: (result: GetCodeReferenceResult) => void, reject): void => {
      client.getCodeReference(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new CodeRefsServiceError(
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

export interface ListCodeRefsResult {
  request: ListCodeReferencesRequest;
  response: ListCodeReferencesResponse;
}

export function listCodeRefs(
  request: ListCodeReferencesRequest
): Promise<ListCodeRefsResult> {
  return new Promise(
    (resolve: (result: ListCodeRefsResult) => void, reject): void => {
      client.listCodeReferences(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new CodeRefsServiceError(
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
