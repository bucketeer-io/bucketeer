import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ArchiveExperimentRequest,
  ArchiveExperimentResponse,
  ArchiveGoalRequest,
  ArchiveGoalResponse,
  CreateExperimentRequest,
  CreateExperimentResponse,
  CreateGoalRequest,
  CreateGoalResponse,
  DeleteGoalRequest,
  DeleteGoalResponse,
  GetExperimentRequest,
  GetExperimentResponse,
  GetGoalRequest,
  GetGoalResponse,
  ListExperimentsRequest,
  ListExperimentsResponse,
  ListGoalsRequest,
  ListGoalsResponse,
  StopExperimentRequest,
  StopExperimentResponse,
  UpdateExperimentRequest,
  UpdateExperimentResponse,
  UpdateGoalRequest,
  UpdateGoalResponse,
} from '../proto/experiment/service_pb';
import {
  ExperimentServiceClient,
  ServiceError,
} from '../proto/experiment/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class ExperimentServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ExperimentServiceError);
    }
    this.name = 'ExperimentServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new ExperimentServiceClient(urls.GRPC);

export interface GetExperimentResult {
  request: GetExperimentRequest;
  response: GetExperimentResponse;
}

export function getExperiment(
  request: GetExperimentRequest
): Promise<GetExperimentResult> {
  return new Promise(
    (resolve: (result: GetExperimentResult) => void, reject): void => {
      client.getExperiment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListExperimentsResult {
  request: ListExperimentsRequest;
  response: ListExperimentsResponse;
}

export function listExperiments(
  request: ListExperimentsRequest
): Promise<ListExperimentsResult> {
  return new Promise(
    (resolve: (result: ListExperimentsResult) => void, reject): void => {
      client.listExperiments(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ExperimentServiceError(
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

export interface CreateExperimentResult {
  request: CreateExperimentRequest;
  response: CreateExperimentResponse;
}

export function createExperiment(
  request: CreateExperimentRequest
): Promise<CreateExperimentResult> {
  return new Promise(
    (resolve: (result: CreateExperimentResult) => void, reject): void => {
      client.createExperiment(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ExperimentServiceError(
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

export interface UpdateExperimentResult {
  request: UpdateExperimentRequest;
  response: UpdateExperimentResponse;
}

export function updateExperiment(
  request: UpdateExperimentRequest
): Promise<UpdateExperimentResult> {
  return new Promise(
    (resolve: (result: UpdateExperimentResult) => void, reject): void => {
      client.updateExperiment(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ExperimentServiceError(
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

export interface ArchiveExperimentResult {
  request: ArchiveExperimentRequest;
  response: ArchiveExperimentResponse;
}

export function archiveExperiment(
  request: ArchiveExperimentRequest
): Promise<ArchiveExperimentResult> {
  return new Promise(
    (resolve: (result: ArchiveExperimentResult) => void, reject): void => {
      client.archiveExperiment(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ExperimentServiceError(
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

export interface StopExperimentResult {
  request: StopExperimentRequest;
  response: StopExperimentResponse;
}

export function stopExperiment(
  request: StopExperimentRequest
): Promise<StopExperimentResult> {
  return new Promise(
    (resolve: (result: StopExperimentResult) => void, reject): void => {
      client.stopExperiment(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface GetGoalResult {
  request: GetGoalRequest;
  response: GetGoalResponse;
}

export function getGoal(request: GetGoalRequest): Promise<GetGoalResult> {
  return new Promise(
    (resolve: (result: GetGoalResult) => void, reject): void => {
      client.getGoal(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListGoalsResult {
  request: ListGoalsRequest;
  response: ListGoalsResponse;
}

export function listGoals(request: ListGoalsRequest): Promise<ListGoalsResult> {
  return new Promise(
    (resolve: (result: ListGoalsResult) => void, reject): void => {
      client.listGoals(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface CreateGoalResult {
  request: CreateGoalRequest;
  response: CreateGoalResponse;
}

export function createGoal(
  request: CreateGoalRequest
): Promise<CreateGoalResult> {
  return new Promise(
    (resolve: (result: CreateGoalResult) => void, reject): void => {
      client.createGoal(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DeleteGoalResult {
  request: DeleteGoalRequest;
  response: DeleteGoalResponse;
}

export function deleteGoal(
  request: DeleteGoalRequest
): Promise<DeleteGoalResult> {
  return new Promise(
    (resolve: (result: DeleteGoalResult) => void, reject): void => {
      client.deleteGoal(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UpdateGoalResult {
  request: UpdateGoalRequest;
  response: UpdateGoalResponse;
}

export function updateGoal(
  request: UpdateGoalRequest
): Promise<UpdateGoalResult> {
  return new Promise(
    (resolve: (result: UpdateGoalResult) => void, reject): void => {
      client.updateGoal(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ArchiveGoalResult {
  request: ArchiveGoalRequest;
  response: ArchiveGoalResponse;
}

export function archiveGoal(
  request: ArchiveGoalRequest
): Promise<ArchiveGoalResult> {
  return new Promise(
    (resolve: (result: ArchiveGoalResult) => void, reject): void => {
      client.archiveGoal(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ExperimentServiceError(
              extractErrorMessage(error),
              request,
              error
            )
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}
