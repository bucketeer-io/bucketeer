import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationDeletePayload {
  id: string;
  environmentId?: string;
}

// DELETE /v1/notification — delete a notification. System admin only.
export const notificationDelete = async (
  payload: NotificationDeletePayload
) => {
  const params = pickBy(payload, v => isNotEmpty(v));

  return axiosClient
    .delete(`/v1/notification?${stringifyParams(params)}`)
    .then(response => response.data);
};
