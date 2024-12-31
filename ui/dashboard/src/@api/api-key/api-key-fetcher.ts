import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { APIKey } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface APIKeyFetcherParams {
  id: string;
  environmentId: string;
}

export interface APIKeyResponse {
  apiKey: Array<APIKey>;
}

export const apiKeyFetcher = async (
  params?: APIKeyFetcherParams
): Promise<APIKeyResponse> => {
  return axiosClient
    .get<APIKeyResponse>('/v1/account/get_api_key', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
