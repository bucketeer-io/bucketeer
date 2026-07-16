import axiosClient from '@api/axios-client';

export interface NotificationMarkAllAsReadPayload {
  environmentId: string;
}

// POST /v1/notifications/mark_all_as_read
export const notificationMarkAllAsRead = async (
  payload: NotificationMarkAllAsReadPayload
) => {
  return axiosClient
    .post('/v1/notifications/mark_all_as_read', payload)
    .then(response => response.data);
};
