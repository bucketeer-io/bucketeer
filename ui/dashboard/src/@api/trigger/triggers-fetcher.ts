import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, TriggerCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface TriggersFetcherParams extends CollectionParams {
  environmentId: string;
  featureId: string;
}

export const triggersFetcher = async (
  _params?: TriggersFetcherParams
): Promise<TriggerCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<TriggerCollection>(`/v1/flag_triggers?${stringifyParams(params)}`)
    .then(response => response.data);
};
