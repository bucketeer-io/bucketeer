import axiosClient from '@api/axios-client';
import { APIKey } from '@types';

export interface APIKeyUpdaterParams {
  id: string;
  environmentId?: string;
  name?: string;
  description?: string;
  disabled?: boolean;
}

export interface APIKeyUpdaterResponse {
  apiKey: Array<APIKey>;
}

export const apiKeyUpdater = async (
  params?: APIKeyUpdaterParams
): Promise<APIKeyUpdaterResponse> => {
  return axiosClient
    .patch<APIKeyUpdaterResponse>('/v1/account/update_api_key', params)
    .then(response => response.data);
};
