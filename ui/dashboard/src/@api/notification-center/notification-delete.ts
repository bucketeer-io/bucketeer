import axiosClient from '@api/axios-client';
import { stringifyParams } from 'utils/search-params';

export interface NotificationDeletePayload {
  id: string;
}

// DELETE /v1/admin_notification — delete a notification. System admin only.
export const notificationDelete = async (
  payload: NotificationDeletePayload
) => {
  return axiosClient
    .delete(`/v1/admin_notification?${stringifyParams(payload)}`)
    .then(response => response.data);
};
