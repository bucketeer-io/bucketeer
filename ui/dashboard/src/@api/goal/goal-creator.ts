import axiosClient from '@api/axios-client';
import { ConnectionType, Goal } from '@types';

export type GoalCreatorPayload = {
  environmentId: string;
  id: string;
  name: string;
  connectionType: ConnectionType;
  description?: string;
};

export interface GoalCreatorResponse {
  goal: Goal;
}

export const goalCreator = async (
  payload?: GoalCreatorPayload
): Promise<GoalCreatorResponse> => {
  return axiosClient
    .post<GoalCreatorResponse>('/v1/goal', payload)
    .then(response => response.data);
};
