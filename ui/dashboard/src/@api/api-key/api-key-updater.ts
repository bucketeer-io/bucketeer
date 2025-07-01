import axiosClient from '@api/axios-client';
import { APIKeyResponse } from './api-key-fetcher';

export interface APIKeyUpdaterPayload {
  id: string;
  environmentId?: string;
  name?: string;
  description?: string;
  disabled?: boolean;
}

export const apiKeyUpdater = async (
  params?: APIKeyUpdaterPayload
): Promise<APIKeyResponse> => {
  return axiosClient
    .patch<APIKeyResponse>('/v1/account/update_api_key', params)
    .then(response => response.data);
};
