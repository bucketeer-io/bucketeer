import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface GoalDeleterParams {
  id: string;
  environmentId: string;
}

export const goalDeleter = async (_params?: GoalDeleterParams) => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .delete(`/v1/goal?${stringifyParams(params)}`)
    .then(response => response.data);
};
