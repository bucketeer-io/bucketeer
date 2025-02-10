import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ArchiveFeatureRequest,
  ArchiveFeatureResponse,
  CreateFeatureRequest,
  CloneFeatureRequest,
  CloneFeatureResponse,
  CreateFeatureResponse,
  DisableFeatureRequest,
  DisableFeatureResponse,
  EnableFeatureRequest,
  EnableFeatureResponse,
  GetFeatureRequest,
  GetFeatureResponse,
  ListFeaturesRequest,
  ListFeaturesResponse,
  UnarchiveFeatureRequest,
  UnarchiveFeatureResponse,
  UpdateFeatureDetailsRequest,
  UpdateFeatureDetailsResponse,
  UpdateFeatureTargetingRequest,
  UpdateFeatureTargetingResponse,
  UpdateFeatureVariationsRequest,
  UpdateFeatureVariationsResponse,
  ListTagsRequest,
  ListTagsResponse,
  EvaluateFeaturesRequest,
  EvaluateFeaturesResponse,
  UpdateFeatureRequest,
  UpdateFeatureResponse
} from '../proto/feature/service_pb';
import {
  FeatureServiceClient,
  ServiceError
} from '../proto/feature/service_pb_service';

import { extractErrorMessage } from './messages';
import {
  checkUnauthenticatedError,
  getMetaDataForClient as getMetaData
} from './utils';
import { UNAUTHENTICATED_ERROR } from '../middlewares/thunkErrorHandler';

export class FeatureServiceError<Request> extends Error {
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
      Error.captureStackTrace(this, FeatureServiceError);
    }
    this.name = 'FeatureServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new FeatureServiceClient(urls.GRPC);

export interface ListFeaturesResult {
  request: ListFeaturesRequest;
  response: ListFeaturesResponse;
}

export function listFeatures(
  request: ListFeaturesRequest
): Promise<ListFeaturesResult> {
  return new Promise(
    (resolve: (result: ListFeaturesResult) => void, reject): void => {
      client.listFeatures(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface CreateFeatureResult {
  request: CreateFeatureRequest;
  response: CreateFeatureResponse;
}

export function createFeature(
  request: CreateFeatureRequest
): Promise<CreateFeatureResult> {
  return new Promise(
    (resolve: (result: CreateFeatureResult) => void, reject): void => {
      client.createFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface CloneFeatureResult {
  request: CloneFeatureRequest;
  response: CloneFeatureResponse;
}

export function cloneFeature(
  request: CloneFeatureRequest
): Promise<CloneFeatureResult> {
  return new Promise(
    (resolve: (result: CloneFeatureResult) => void, reject): void => {
      client.cloneFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface EnableFeatureResult {
  request: EnableFeatureRequest;
  response: EnableFeatureResponse;
}

export function enableFeature(
  request: EnableFeatureRequest
): Promise<EnableFeatureResult> {
  return new Promise(
    (resolve: (result: EnableFeatureResult) => void, reject): void => {
      client.enableFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DisableFeatureResult {
  request: DisableFeatureRequest;
  response: DisableFeatureResponse;
}

export function disableFeature(
  request: DisableFeatureRequest
): Promise<DisableFeatureResult> {
  return new Promise(
    (resolve: (result: DisableFeatureResult) => void, reject): void => {
      client.disableFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ArchiveFeatureResult {
  request: ArchiveFeatureRequest;
  response: ArchiveFeatureResponse;
}

export function archiveFeature(
  request: ArchiveFeatureRequest
): Promise<ArchiveFeatureResult> {
  return new Promise(
    (resolve: (result: ArchiveFeatureResult) => void, reject): void => {
      client.archiveFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UnarchiveFeatureResult {
  request: UnarchiveFeatureRequest;
  response: UnarchiveFeatureResponse;
}

export function unarchiveFeature(
  request: UnarchiveFeatureRequest
): Promise<UnarchiveFeatureResult> {
  return new Promise(
    (resolve: (result: UnarchiveFeatureResult) => void, reject): void => {
      client.unarchiveFeature(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FeatureServiceError(
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

export interface GetFeatureResult {
  request: GetFeatureRequest;
  response: GetFeatureResponse;
}

export function getFeature(
  request: GetFeatureRequest
): Promise<GetFeatureResult> {
  return new Promise(
    (resolve: (result: GetFeatureResult) => void, reject): void => {
      client.getFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

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
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UpdateFeatureResult {
  request: UpdateFeatureRequest;
  response: UpdateFeatureResponse;
}

export function updateFeature(
  request: UpdateFeatureRequest
): Promise<UpdateFeatureResult> {
  return new Promise(
    (resolve: (result: UpdateFeatureResult) => void, reject): void => {
      client.updateFeature(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new FeatureServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UpdateFeatureTargetingResult {
  request: UpdateFeatureTargetingRequest;
  response: UpdateFeatureTargetingResponse;
}

export function updateFeatureTargeting(
  request: UpdateFeatureTargetingRequest
): Promise<UpdateFeatureTargetingResult> {
  return new Promise(
    (resolve: (result: UpdateFeatureTargetingResult) => void, reject): void => {
      client.updateFeatureTargeting(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FeatureServiceError(
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

export interface UpdateFeatureVariationsResult {
  request: UpdateFeatureVariationsRequest;
  response: UpdateFeatureVariationsResponse;
}

export function updateFeatureVariations(
  request: UpdateFeatureVariationsRequest
): Promise<UpdateFeatureVariationsResult> {
  return new Promise(
    (
      resolve: (result: UpdateFeatureVariationsResult) => void,
      reject
    ): void => {
      client.updateFeatureVariations(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FeatureServiceError(
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

export interface UpdateFeatureDetailsResult {
  request: UpdateFeatureDetailsRequest;
  response: UpdateFeatureDetailsResponse;
}

export function updateFeatureDetails(
  request: UpdateFeatureDetailsRequest
): Promise<UpdateFeatureDetailsResult> {
  return new Promise(
    (resolve: (result: UpdateFeatureDetailsResult) => void, reject): void => {
      client.updateFeatureDetails(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FeatureServiceError(
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

export interface EvaluateFeaturesResult {
  request: EvaluateFeaturesRequest;
  response: EvaluateFeaturesResponse;
}

export function evaluateFeatures(
  request: EvaluateFeaturesRequest
): Promise<EvaluateFeaturesResult> {
  return new Promise(
    (resolve: (result: EvaluateFeaturesResult) => void, reject): void => {
      client.evaluateFeatures(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new FeatureServiceError(
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
