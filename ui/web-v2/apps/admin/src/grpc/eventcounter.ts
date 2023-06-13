import { grpc } from '@improbable-eng/grpc-web';
import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import { intl } from '../lang';
import { messages } from '../lang/messages';
import {
  GetEvaluationTimeseriesCountRequest,
  GetEvaluationTimeseriesCountResponse,
  GetExperimentResultRequest,
  GetExperimentResultResponse,
} from '../proto/eventcounter/service_pb';
import {
  EventCounterServiceClient,
  ServiceError,
} from '../proto/eventcounter/service_pb_service';

import { deserializeStatus, extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class EventCounterServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, EventCounterServiceError);
    }
    this.name = 'EventCounterServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new EventCounterServiceClient(urls.GRPC);

export interface GetExperimentResultResult {
  request: GetExperimentResultRequest;
  response: GetExperimentResultResponse;
}

export function getExperimentResult(
  request: GetExperimentResultRequest
): Promise<GetExperimentResultResult> {
  return new Promise(
    (resolve: (result: GetExperimentResultResult) => void, reject): void => {
      client.getExperimentResult(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            // Not found error returns only if experiment started,
            // but the resutl is not yet created.
            if (
              deserializeStatus(error.metadata).getCode() === grpc.Code.NotFound
            ) {
              reject(
                new EventCounterServiceError(
                  intl.formatMessage(
                    messages.experiment.result.noData.errorMessage
                  ),
                  request,
                  error
                )
              );
              return;
            }
            reject(
              new EventCounterServiceError(
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

export interface GetEvaluationTimeseriesCountResult {
  request: GetEvaluationTimeseriesCountRequest;
  response: GetEvaluationTimeseriesCountResponse;
}

export function getEvaluationTimeseriesCount(
  request: GetEvaluationTimeseriesCountRequest
): Promise<GetEvaluationTimeseriesCountResult> {
  return new Promise(
    (
      resolve: (result: GetEvaluationTimeseriesCountResult) => void,
      reject
    ): void => {
      client.getEvaluationTimeseriesCount(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new EventCounterServiceError(
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
