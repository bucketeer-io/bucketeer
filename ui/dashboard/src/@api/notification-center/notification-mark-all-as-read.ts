import axiosClient from '@api/axios-client';

// POST /v1/notifications/mark_all_as_read
export const notificationMarkAllAsRead = async () => {
  return axiosClient
    .post('/v1/notifications/mark_all_as_read', {})
    .then(response => response.data);
};
