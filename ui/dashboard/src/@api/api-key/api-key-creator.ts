import axiosClient from '@api/axios-client';
import { APIKey, APIKeyRole } from '@types';

export interface APIKeyCreatorParams {
  environmentId: string;
  name: string;
  role: APIKeyRole;
  description?: string;
}

export interface APIKeyCreatorResponse {
  apiKey: APIKey;
}

export const apiKeyCreator = async (
  params?: APIKeyCreatorParams
): Promise<APIKeyCreatorResponse> => {
  return axiosClient
    .post<APIKeyCreatorResponse>('/v1/account/create_api_key', params)
    .then(response => response.data);
};
