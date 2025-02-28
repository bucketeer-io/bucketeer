import axiosClient from '@api/axios-client';

export interface EnvironmentArchiveParams {
  id: string;
}

export const environmentArchive = async (params?: EnvironmentArchiveParams) => {
  return axiosClient
    .post('/v1/environment/archive_environment', params)
    .then(response => response.data);
};
