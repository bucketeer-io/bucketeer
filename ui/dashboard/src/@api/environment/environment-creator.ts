import axiosClient from '@api/axios-client';
import { Environment } from '@types';

export interface EnvironmentCreatorParams {
  name: string;
  urlCode: string;
  description?: string;
  projectId: string;
  requireComment: boolean;
}

export interface EnvironmentResponse {
  environment: Environment;
}

export const environmentCreator = async (
  params?: EnvironmentCreatorParams
): Promise<EnvironmentResponse> => {
  return axiosClient
    .post<EnvironmentResponse>('/v1/environment/create_environment', params)
    .then(response => response.data);
};
