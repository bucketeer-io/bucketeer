import axiosClient from '@api/axios-client';

export interface NotificationMarkAsReadPayload {
  environmentId: string;
  ids: string[];
}

// POST /v1/notifications/mark_as_read — idempotent.
export const notificationMarkAsRead = async (
  payload: NotificationMarkAsReadPayload
) => {
  return axiosClient
    .post('/v1/notifications/mark_as_read', payload)
    .then(response => response.data);
};
