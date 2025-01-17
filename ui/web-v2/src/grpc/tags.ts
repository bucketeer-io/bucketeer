import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';

import { extractErrorMessage } from './messages';
import {
  checkUnauthenticatedError,
  getMetaDataForClient as getMetaData
} from './utils';
import { UNAUTHENTICATED_ERROR } from '../middlewares/thunkErrorHandler';
import {
  TagServiceClient,
  ServiceError
} from '../proto/tag/service_pb_service';
import { ListTagsRequest, ListTagsResponse } from '../proto/tag/service_pb';

const client = new TagServiceClient(urls.GRPC);

export interface ListTagsResult {
  request: ListTagsRequest;
  response: ListTagsResponse;
}

export function listTags(request: ListTagsRequest): Promise<ListTagsResult> {
  return new Promise(
    (resolve: (result: ListTagsResult) => void, reject): void => {
      client.listTags(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new TagServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export class TagServiceError<Request> extends Error {
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
      Error.captureStackTrace(this, TagServiceError);
    }
    this.name = 'TagServiceError';
    this.request = request;
    this.error = error;
  }
}
