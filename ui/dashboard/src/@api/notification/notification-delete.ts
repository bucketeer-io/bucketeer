import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationDeletePayload {
  id: string;
  environmentId?: string;
}

export const notificationDelete = async (
  payload: NotificationDeletePayload
) => {
  const _payload = pickBy(payload, v => isNotEmpty(v));

  return axiosClient
    .delete(`/v1/subscription?${stringifyParams(_payload)}`)
    .then(response => response.data);
};
