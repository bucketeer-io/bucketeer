import axiosClient from '@api/axios-client';
import type { ExecuteScheduledFlagChangeResponse } from '@types';

export interface ScheduledFlagChangeExecuteParams {
  environmentId: string;
  id: string;
}

export const scheduledFlagChangeExecutor = async (
  params?: ScheduledFlagChangeExecuteParams
): Promise<ExecuteScheduledFlagChangeResponse> => {
  return axiosClient
    .post<ExecuteScheduledFlagChangeResponse>(
      '/v1/scheduled_flag_change/execute',
      params
    )
    .then(response => response.data);
};
