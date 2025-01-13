import axiosClient from '@api/axios-client';
import { Goal } from '@types';

export type GoalCreatorPayload = {
  environmentId: string;
  id: string;
  name: string;
  description?: string;
};

export interface GoalCreatorResponse {
  goal: Array<Goal>;
}

export const goalCreator = async (
  payload?: GoalCreatorPayload
): Promise<GoalCreatorResponse> => {
  return axiosClient
    .post<GoalCreatorResponse>('/v1/goal', payload)
    .then(response => response.data);
};
