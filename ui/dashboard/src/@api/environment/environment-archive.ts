import axiosClient from '@api/axios-client';
import { EnvironmentResponse } from './environment-creator';

export interface EnvironmentArchiveParams {
  id: string;
}

export const environmentArchive = async (
  params?: EnvironmentArchiveParams
): Promise<EnvironmentResponse> => {
  return axiosClient
    .post<EnvironmentResponse>('/v1/environment/update_environment', {
      id: params?.id,
      archived: true
    })
    .then(response => response.data);
};
