import axiosClient from '@api/axios-client';
import { APIKey, APIKeyRole } from '@types';

export interface APIKeyCreatorCommand {
  name: string;
  role: APIKeyRole;
}

export interface APIKeyCreatorParams {
  environmentId: string;
  command: APIKeyCreatorCommand;
}

export interface APIKeyCreatorResponse {
  apiKey: Array<APIKey>;
}

export const apiKeyCreator = async (
  params?: APIKeyCreatorParams
): Promise<APIKeyCreatorResponse> => {
  return axiosClient
    .post<APIKeyCreatorResponse>('/v1/account/create_api_key', params)
    .then(response => response.data);
};
