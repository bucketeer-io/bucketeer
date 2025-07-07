import axiosClient from '@api/axios-client';
import { EnvironmentResponse } from './environment-creator';

export interface EnvironmentUpdateParams {
  id: string;
  name?: string;
  description?: string;
  requireComment?: boolean;
  archived?: boolean;
}

export const environmentUpdater = async (
  params?: EnvironmentUpdateParams
): Promise<EnvironmentResponse> => {
  return axiosClient
    .post<EnvironmentResponse>('/v1/environment/update_environment', params)
    .then(response => response.data);
};
