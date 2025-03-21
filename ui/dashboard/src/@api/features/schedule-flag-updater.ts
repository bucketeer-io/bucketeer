import axiosClient from '@api/axios-client';
import { ScheduleChangeItem } from '@types';

export interface ScheduleUpdaterParams {
  id: string;
  environmentId: string;
  scheduledAt: string;
  scheduledChanges: ScheduleChangeItem[];
}

export const scheduleFlagUpdater = async (params?: ScheduleUpdaterParams) => {
  return axiosClient
    .patch('/v1/schedule_flag', params)
    .then(response => response.data);
};
