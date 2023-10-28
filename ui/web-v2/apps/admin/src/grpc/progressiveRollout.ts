import { Nullable, isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { urls } from '../config';
import {
  CreateProgressiveRolloutRequest,
  CreateProgressiveRolloutResponse,
  DeleteProgressiveRolloutRequest,
  DeleteProgressiveRolloutResponse,
  ListProgressiveRolloutsRequest,
  ListProgressiveRolloutsResponse,
  // ExecuteProgressiveRolloutRequest,
  // ExecuteProgressiveRolloutResponse,
} from '../proto/autoops/service_pb';
import {
  AutoOpsServiceClient,
  ServiceError,
} from '../proto/autoops/service_pb_service';

import { extractErrorMessage } from './messages';
import { getMetaDataForClient as getMetaData } from './utils';

export class ProgressiveRolloutServiceError<Request> extends Error {
  request: Request;

  error: Nullable<ServiceError>;

  constructor(
    message: string,
    request: Request,
    error: Nullable<ServiceError>
  ) {
    super(message);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ProgressiveRolloutServiceError);
    }
    this.name = 'ProgressiveRolloutServiceError';
    this.request = request;
    this.error = error;
  }
}

const client = new AutoOpsServiceClient(urls.GRPC);

export interface CreateProgressiveRolloutResult {
  request: CreateProgressiveRolloutRequest;
  response: CreateProgressiveRolloutResponse;
}

export function createProgressiveRollout(
  request: CreateProgressiveRolloutRequest
): Promise<CreateProgressiveRolloutResult> {
  return new Promise(
    (
      resolve: (result: CreateProgressiveRolloutResult) => void,
      reject
    ): void => {
      client.createProgressiveRollout(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ProgressiveRolloutServiceError(
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

export interface ListProgressiveRolloutsResult {
  request: ListProgressiveRolloutsRequest;
  response: ListProgressiveRolloutsResponse;
}

export function listProgressiveRollouts(
  request: ListProgressiveRolloutsRequest
): Promise<ListProgressiveRolloutsResult> {
  return new Promise(
    (
      resolve: (result: ListProgressiveRolloutsResult) => void,
      reject
    ): void => {
      client.listProgressiveRollouts(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ProgressiveRolloutServiceError(
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

// export interface UpdateAutoOpsRuleResult {
//   request: UpdateAutoOpsRuleRequest;
//   response: UpdateAutoOpsRuleResponse;
// }

// export function updateAutoOpsRule(
//   request: UpdateAutoOpsRuleRequest
// ): Promise<UpdateAutoOpsRuleResult> {
//   return new Promise(
//     (resolve: (result: UpdateAutoOpsRuleResult) => void, reject): void => {
//       client.updateAutoOpsRule(
//         request,
//         getMetaData(),
//         (error, response): void => {
//           if (isNotNull(error) || isNull(response)) {
//             reject(
//               new ProgressiveRolloutServiceError(
//                 extractErrorMessage(error),
//                 request,
//                 error
//               )
//             );
//           } else {
//             resolve({ request, response });
//           }
//         }
//       );
//     }
//   );
// }

export interface DeleteProgressiveRolloutResult {
  request: DeleteProgressiveRolloutRequest;
  response: DeleteProgressiveRolloutResponse;
}

export function deleteProgressiveRollout(
  request: DeleteProgressiveRolloutRequest
): Promise<DeleteProgressiveRolloutResult> {
  return new Promise(
    (
      resolve: (result: DeleteProgressiveRolloutResult) => void,
      reject
    ): void => {
      client.deleteProgressiveRollout(
        request,
        getMetaData(),
        (error, response): void => {
          if (isNotNull(error) || isNull(response)) {
            reject(
              new ProgressiveRolloutServiceError(
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

// export interface ListOpsCountsResult {
//   request: ListOpsCountsRequest;
//   response: ListOpsCountsResponse;
// }

// export function listOpsCounts(
//   request: ListOpsCountsRequest
// ): Promise<ListOpsCountsResult> {
//   return new Promise(
//     (resolve: (result: ListOpsCountsResult) => void, reject): void => {
//       client.listOpsCounts(request, getMetaData(), (error, response): void => {
//         if (isNotNull(error) || isNull(response)) {
//           reject(
//             new ProgressiveRolloutServiceError(extractErrorMessage(error), request, error)
//           );
//         } else {
//           resolve({ request, response });
//         }
//       });
//     }
//   );
// }
