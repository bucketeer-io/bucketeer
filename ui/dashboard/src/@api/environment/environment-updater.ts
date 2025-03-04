import axiosClient from '@api/axios-client';

export interface EnvironmentUpdateParams {
  id: string;
  name: string;
  description?: string;
  requireComment: boolean;
}

export const environmentUpdater = async (params?: EnvironmentUpdateParams) => {
  return axiosClient
    .post('/v1/environment/update_environment', params)
    .then(response => response.data);
};
