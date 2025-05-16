import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { Rollout } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface RolloutFetcherParams {
  environmentId: string;
  id: string;
}

export interface RolloutFetcherResponse {
  progressiveRollout: Rollout;
}

export const rolloutFetcher = async (
  params?: RolloutFetcherParams
): Promise<RolloutFetcherResponse> => {
  const _params = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<RolloutFetcherResponse>(`/v1/progressive_rollout?${_params}`)
    .then(response => response.data);
};
