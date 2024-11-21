import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { APIKeyCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface APIKeyFetcherParams extends CollectionParams {
  environmentNamespace?: string;
}

export const apiKeysFetcher = async (
  params?: APIKeyFetcherParams
): Promise<APIKeyCollection> => {
  return axiosClient
    .post<APIKeyCollection>(
      '/v1/account/list_api_keys',
      pickBy(params, v => isNotEmpty(v))
    )
    .then(response => response.data);
};
