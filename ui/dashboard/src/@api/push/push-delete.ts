import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface PushDeletePayload {
  id: string;
  environmentId?: string;
}

export const pushDelete = async (payload: PushDeletePayload) => {
  const _payload = pickBy(payload, v => isNotEmpty(v));

  return axiosClient
    .delete(`/v1/push?${stringifyParams(_payload)}`)
    .then(response => response.data);
};
