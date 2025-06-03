import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AutoOpsDeleteParams {
  id: string;
  environmentId: string;
}

export const autoOpsDelete = async (_params?: AutoOpsDeleteParams) => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .delete(`/v1/auto_ops_rule?${stringifyParams(params)}`)
    .then(response => response.data);
};
