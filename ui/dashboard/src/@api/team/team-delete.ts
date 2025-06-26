import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface TeamDeletePayload {
  id: string;
  organizationId: string;
}

export const teamDelete = async (params?: TeamDeletePayload) => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .delete(`/v1/team?${requestParams}`)
    .then(response => response.data);
};
