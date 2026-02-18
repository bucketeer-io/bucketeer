import axiosClient from '@api/axios-client';
import type {
  ScheduledChangePayload,
  UpdateScheduledFlagChangeResponse
} from '@types';

export interface ScheduledFlagChangeUpdaterParams {
  environmentId: string;
  id: string;
  scheduledAt?: string;
  timezone?: string;
  payload?: ScheduledChangePayload;
  comment?: string;
  ignoreConflicts?: boolean;
}

export const scheduledFlagChangeUpdater = async (
  params?: ScheduledFlagChangeUpdaterParams
): Promise<UpdateScheduledFlagChangeResponse> => {
  return axiosClient
    .patch<UpdateScheduledFlagChangeResponse>(
      '/v1/scheduled_flag_change',
      params
    )
    .then(response => response.data);
};
