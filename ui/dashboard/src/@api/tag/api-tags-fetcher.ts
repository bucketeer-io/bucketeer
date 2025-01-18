import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { TagCollection, CollectionParams, EntityType } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface TagsFetcherParams extends CollectionParams {
  environmentId?: string;
  entityType?: EntityType;
}

export const tagsFetcher = async (
  params?: TagsFetcherParams
): Promise<TagCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<TagCollection>(`/v1/tags?${requestParams}`)
    .then(response => response.data);
};
