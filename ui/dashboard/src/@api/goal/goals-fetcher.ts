import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, ConnectionType, GoalCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface GoalsFetcherParams extends CollectionParams {
  environmentId?: string;
  isInUseStatus?: boolean;
  archived?: boolean;
  connectionType?: ConnectionType;
}

export const goalsFetcher = async (
  params?: GoalsFetcherParams
): Promise<GoalCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<GoalCollection>(`/v1/goals?${requestParams}`)
    .then(response => response.data);
};
