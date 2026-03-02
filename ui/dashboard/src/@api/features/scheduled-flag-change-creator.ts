import axiosClient from '@api/axios-client';
import type {
  ScheduledChangePayload,
  CreateScheduledFlagChangeResponse
} from '@types';

export interface ScheduledFlagChangeCreatorParams {
  environmentId: string;
  featureId: string;
  scheduledAt: string;
  timezone?: string;
  payload: ScheduledChangePayload;
  comment?: string;
  ignoreConflicts?: boolean;
}

export const scheduledFlagChangeCreator = async (
  params?: ScheduledFlagChangeCreatorParams
): Promise<CreateScheduledFlagChangeResponse> => {
  return axiosClient
    .post<CreateScheduledFlagChangeResponse>(
      '/v1/scheduled_flag_change',
      params
    )
    .then(response => response.data);
};
