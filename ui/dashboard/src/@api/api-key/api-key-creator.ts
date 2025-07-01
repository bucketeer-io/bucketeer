import axiosClient from '@api/axios-client';
import { APIKeyRole } from '@types';
import { APIKeyResponse } from './api-key-fetcher';

export interface APIKeyCreatorPayload {
  environmentId: string;
  name: string;
  role: APIKeyRole;
  description?: string;
}

export const apiKeyCreator = async (
  params?: APIKeyCreatorPayload
): Promise<APIKeyResponse> => {
  return axiosClient
    .post<APIKeyResponse>('/v1/account/create_api_key', params)
    .then(response => response.data);
};
