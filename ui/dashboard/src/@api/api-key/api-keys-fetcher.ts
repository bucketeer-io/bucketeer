import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { APIKeyCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface APIKeysFetcherParams extends CollectionParams {
  environmentId?: string;
  environmentIds?: string[];
  organizationId?: string;
}

export const apiKeysFetcher = async (
  params?: APIKeysFetcherParams
): Promise<APIKeyCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<APIKeyCollection>(`/v1/account/list_api_keys?${requestParams}`)
    .then(response => response.data);
};
