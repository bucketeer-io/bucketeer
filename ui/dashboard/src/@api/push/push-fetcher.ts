import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { Push } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface PushFetcherParams {
  id: string;
  environmentId: string;
}

export interface PushResponse {
  push: Array<Push>;
}

export const pushFetcher = async (
  params?: PushFetcherParams
): Promise<PushResponse> => {
  return axiosClient
    .get<PushResponse>('/v1/push', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
