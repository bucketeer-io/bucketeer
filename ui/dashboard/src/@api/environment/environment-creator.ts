import axiosClient from '@api/axios-client';
import { Environment } from '@types';

export interface EnvironmentCreatorCommand {
  name: string;
  urlCode: string;
  description?: string;
  projectId: string;
  requireComment: boolean;
}

export interface EnvironmentCreatorParams {
  command: EnvironmentCreatorCommand;
}

export interface EnvironmentResponse {
  environment: Array<Environment>;
}

export const environmentCreator = async (
  params?: EnvironmentCreatorParams
): Promise<EnvironmentResponse> => {
  return axiosClient
    .post<EnvironmentResponse>('/v1/environment/create_environment', params)
    .then(response => response.data);
};
