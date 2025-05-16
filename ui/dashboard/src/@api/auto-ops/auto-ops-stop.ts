import axiosClient from '@api/axios-client';

export interface AutoOpsStopParams {
  id: string;
  environmentId: string;
}

export const autoOpsStop = async (params?: AutoOpsStopParams) => {
  return axiosClient
    .post('/v1/auto_ops_rule/stop', params)
    .then(response => response.data);
};
