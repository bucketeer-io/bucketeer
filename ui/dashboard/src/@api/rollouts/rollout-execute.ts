import axiosClient from '@api/axios-client';

export interface RolloutExecuteParams {
  id: string;
  environmentId: string;
  scheduleId: string;
}

export const rolloutExecute = async (params?: RolloutExecuteParams) => {
  return axiosClient
    .post('/v1/progressive_rollout/execute', params)
    .then(response => response.data);
};
