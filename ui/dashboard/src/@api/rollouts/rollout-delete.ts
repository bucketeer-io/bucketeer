import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface RolloutDeleteParams {
  environmentId: string;
  id: string;
}

export const rolloutDelete = async (params?: RolloutDeleteParams) => {
  const _params = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .delete(`/v1/progressive_rollout?${_params}`)
    .then(response => response.data);
};
