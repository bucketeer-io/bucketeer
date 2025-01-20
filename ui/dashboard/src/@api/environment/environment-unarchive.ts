import axiosClient from '@api/axios-client';

export interface EnvironmentUnarchiveParams {
  id: string;
}

export const environmentUnarchive = async (
  params?: EnvironmentUnarchiveParams
) => {
  return axiosClient
    .post('/v1/environment/unarchive_environment', params)
    .then(response => response.data);
};
