import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  AddSegmentUserRequest,
  AddSegmentUserResponse,
  BulkDownloadSegmentUsersRequest,
  BulkDownloadSegmentUsersResponse,
  BulkUploadSegmentUsersRequest,
  BulkUploadSegmentUsersResponse,
  CreateSegmentRequest,
  CreateSegmentResponse,
  DeleteSegmentRequest,
  DeleteSegmentResponse,
  DeleteSegmentUserRequest,
  DeleteSegmentUserResponse,
  GetSegmentRequest,
  GetSegmentResponse,
  ListSegmentUsersRequest,
  ListSegmentUsersResponse,
  ListSegmentsRequest,
  ListSegmentsResponse,
  UpdateSegmentRequest,
  UpdateSegmentResponse,
} from '../proto/feature/service_pb';
import {
  FeatureServiceClient,
  ServiceError,
} from '../proto/feature/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class SegmentServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, SegmentServiceError);
    }
    this.name = 'SegmentServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new FeatureServiceClient(urls.GRPC);

export interface CreateSegmentResult {
  request: CreateSegmentRequest;
  response: CreateSegmentResponse;
}

export function createSegment(
  request: CreateSegmentRequest
): Promise<CreateSegmentResult> {
  return new Promise(
    (resolve: (result: CreateSegmentResult) => void, reject): void => {
      client.createSegment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new SegmentServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListSegmentsResult {
  request: ListSegmentsRequest;
  response: ListSegmentsResponse;
}

export function listSegments(
  request: ListSegmentsRequest
): Promise<ListSegmentsResult> {
  return new Promise(
    (resolve: (result: ListSegmentsResult) => void, reject): void => {
      client.listSegments(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new SegmentServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UpdateSegmentResult {
  request: UpdateSegmentRequest;
  response: UpdateSegmentResponse;
}

export function updateSegment(
  request: UpdateSegmentRequest
): Promise<UpdateSegmentResult> {
  return new Promise(
    (resolve: (result: UpdateSegmentResult) => void, reject): void => {
      client.updateSegment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new SegmentServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DeleteSegmentResult {
  request: DeleteSegmentRequest;
  response: DeleteSegmentResponse;
}

export function deleteSegment(
  request: DeleteSegmentRequest
): Promise<DeleteSegmentResult> {
  return new Promise(
    (resolve: (result: DeleteSegmentResult) => void, reject): void => {
      client.deleteSegment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new SegmentServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface GetSegmentResult {
  request: GetSegmentRequest;
  response: GetSegmentResponse;
}

export function getSegment(
  request: GetSegmentRequest
): Promise<GetSegmentResult> {
  return new Promise(
    (resolve: (result: GetSegmentResult) => void, reject): void => {
      client.getSegment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new SegmentServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListSegmentUsersResult {
  request: ListSegmentUsersRequest;
  response: ListSegmentUsersResponse;
}

export function listSegmentUsers(
  request: ListSegmentUsersRequest
): Promise<ListSegmentUsersResult> {
  return new Promise(
    (resolve: (result: ListSegmentUsersResult) => void, reject): void => {
      client.listSegmentUsers(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new SegmentServiceError(
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

export interface AddSegmentUserResult {
  request: AddSegmentUserRequest;
  response: AddSegmentUserResponse;
}

export function addSegmentUser(
  request: AddSegmentUserRequest
): Promise<AddSegmentUserResult> {
  return new Promise(
    (resolve: (result: AddSegmentUserResult) => void, reject): void => {
      client.addSegmentUser(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new SegmentServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DeleteSegmentUserResult {
  request: DeleteSegmentUserRequest;
  response: DeleteSegmentUserResponse;
}

export function deleteSegmentUser(
  request: DeleteSegmentUserRequest
): Promise<DeleteSegmentUserResult> {
  return new Promise(
    (resolve: (result: DeleteSegmentUserResult) => void, reject): void => {
      client.deleteSegmentUser(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new SegmentServiceError(
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

export interface BulkUploadSegmentUsersResult {
  request: BulkUploadSegmentUsersRequest;
  response: BulkUploadSegmentUsersResponse;
}

export function bulkUploadSegmentUsers(
  request: BulkUploadSegmentUsersRequest
): Promise<BulkUploadSegmentUsersResult> {
  return new Promise(
    (resolve: (result: BulkUploadSegmentUsersResult) => void, reject): void => {
      client.bulkUploadSegmentUsers(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new SegmentServiceError(
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

export interface BulkDownloadSegmentUsersResult {
  request: BulkDownloadSegmentUsersRequest;
  response: BulkDownloadSegmentUsersResponse;
}

export function bulkDownloadSegmentUsers(
  request: BulkDownloadSegmentUsersRequest
): Promise<BulkDownloadSegmentUsersResult> {
  return new Promise(
    (
      resolve: (result: BulkDownloadSegmentUsersResult) => void,
      reject
    ): void => {
      client.bulkDownloadSegmentUsers(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new SegmentServiceError(
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
