import axiosClient from '@api/axios-client';
import { EnvironmentResponse } from './environment-creator';

export interface EnvironmentUnarchiveParams {
  id: string;
}

export const environmentUnarchive = async (
  params?: EnvironmentUnarchiveParams
): Promise<EnvironmentResponse> => {
  return axiosClient
    .post<EnvironmentResponse>('/v1/environment/update_environment', {
      id: params?.id,
      archived: false
    })
    .then(response => response.data);
};
