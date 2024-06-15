import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  ConvertTrialProjectRequest,
  ConvertTrialProjectResponse,
  CreateProjectRequest,
  CreateProjectResponse,
  DisableProjectRequest,
  DisableProjectResponse,
  EnableProjectRequest,
  EnableProjectResponse,
  GetProjectRequest,
  GetProjectResponse,
  ListProjectsRequest,
  ListProjectsResponse,
  UpdateProjectRequest,
  UpdateProjectResponse,
} from '../proto/environment/service_pb';
import {
  EnvironmentServiceClient,
  ServiceError,
} from '../proto/environment/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class ProjectServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ProjectServiceError);
    }
    this.name = 'ProjectServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new EnvironmentServiceClient(urls.GRPC);

export interface GetProjectResult {
  request: GetProjectRequest;
  response: GetProjectResponse;
}

export function getProject(
  request: GetProjectRequest
): Promise<GetProjectResult> {
  return new Promise(
    (resolve: (result: GetProjectResult) => void, reject): void => {
      client.getProject(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ProjectServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ListProjectsResult {
  request: ListProjectsRequest;
  response: ListProjectsResponse;
}

export function listProjects(
  request: ListProjectsRequest
): Promise<ListProjectsResult> {
  return new Promise(
    (resolve: (result: ListProjectsResult) => void, reject): void => {
      client.listProjects(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ProjectServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface CreateProjectResult {
  request: CreateProjectRequest;
  response: CreateProjectResponse;
}

export function createProject(
  request: CreateProjectRequest
): Promise<CreateProjectResult> {
  return new Promise(
    (resolve: (result: CreateProjectResult) => void, reject): void => {
      client.createProject(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ProjectServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface UpdateProjectResult {
  request: UpdateProjectRequest;
  response: UpdateProjectResponse;
}

export function updateProject(
  request: UpdateProjectRequest
): Promise<UpdateProjectResult> {
  return new Promise(
    (resolve: (result: UpdateProjectResult) => void, reject): void => {
      client.updateProject(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ProjectServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface EnableProjectResult {
  request: EnableProjectRequest;
  response: EnableProjectResponse;
}

export function enableProject(
  request: EnableProjectRequest
): Promise<EnableProjectResult> {
  return new Promise(
    (resolve: (result: EnableProjectResult) => void, reject): void => {
      client.enableProject(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ProjectServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface DisableProjectResult {
  request: DisableProjectRequest;
  response: DisableProjectResponse;
}

export function disableProject(
  request: DisableProjectRequest
): Promise<DisableProjectResult> {
  return new Promise(
    (resolve: (result: DisableProjectResult) => void, reject): void => {
      client.disableProject(request, getMetaData(), (error, response): void => {
        if (isNotNull(error) || isNull(response)) {
          reject(
            new ProjectServiceError(extractErrorMessage(error), request, error)
          );
        } else {
          resolve({ request, response });
        }
      });
    }
  );
}

export interface ConvertTrialProjectResult {
  request: ConvertTrialProjectRequest;
  response: ConvertTrialProjectResponse;
}

export function convertTrialProject(
  request: ConvertTrialProjectRequest
): Promise<ConvertTrialProjectResult> {
  return new Promise(
    (resolve: (result: ConvertTrialProjectResult) => void, reject): void => {
      client.convertTrialProject(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ProjectServiceError(
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
