import axiosClient from '@api/axios-client';
import { ScheduleChangeItem } from '@types';

export interface ScheduleCreatorParams {
  featureId: string;
  environmentId: string;
  scheduledAt: string;
  scheduledChanges: ScheduleChangeItem[];
}

export const scheduleFlagCreator = async (params?: ScheduleCreatorParams) => {
  return axiosClient
    .post('/v1/schedule_flag', params)
    .then(response => response.data);
};
