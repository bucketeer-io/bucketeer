import axiosClient from '@api/axios-client';
import { GoalCreatorResponse } from './goal-creator';

export type GoalUpdaterPayload = {
  id: string;
  name: string;
  environmentId: string;
  description?: string;
  archived?: boolean;
};

export const goalUpdater = async (
  payload?: GoalUpdaterPayload
): Promise<GoalCreatorResponse> => {
  return axiosClient
    .patch<GoalCreatorResponse>('/v1/goal', payload)
    .then(response => response.data);
};
