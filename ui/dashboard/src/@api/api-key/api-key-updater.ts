import axiosClient from '@api/axios-client';
import { APIKey } from '@types';

export interface APIKeyUpdaterCommand {
  name: string;
}

export interface APIKeyUpdaterParams {
  id: string;
  environmentId: string;
  command: APIKeyUpdaterCommand;
}

export interface APIKeyUpdaterResponse {
  apiKey: Array<APIKey>;
}

export const apiKeyUpdater = async (
  params?: APIKeyUpdaterParams
): Promise<APIKeyUpdaterResponse> => {
  return axiosClient
    .post<APIKeyUpdaterResponse>('/v1/account/update_api_key', params)
    .then(response => response.data);
};
