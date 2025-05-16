import axiosClient from '@api/axios-client';
import { RolloutStoppedBy } from '@types';

export interface RolloutStoppedParams {
  id: string;
  environmentId: string;
  stoppedBy: RolloutStoppedBy;
}

export const rolloutStopped = async (params?: RolloutStoppedParams) => {
  return axiosClient
    .patch('/v1/progressive_rollout/stop', params)
    .then(response => response.data);
};
