import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface UserSegmentDeleteParams {
  id: string;
  environmentId: string;
}

export const userSegmentDelete = async (_params?: UserSegmentDeleteParams) => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .delete(`/v1/segment?${stringifyParams(params)}`)
    .then(response => response.data);
};
