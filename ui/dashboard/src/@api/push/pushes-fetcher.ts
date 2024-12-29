import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { PushCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface PushesFetcherParams extends CollectionParams {
  environmentId?: string;
}

export const pushesFetcher = async (
  params?: PushesFetcherParams
): Promise<PushCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<PushCollection>(`/v1/pushes?${requestParams}`)
    .then(response => response.data);
};
