import axiosClient from '@api/axios-client';

interface UnreadCountWire {
  count: string;
}

// GET /v1/notifications/unread_count — unread count for the bell badge.
export const notificationUnreadCountFetcher = async (): Promise<number> => {
  return axiosClient
    .get<UnreadCountWire>('/v1/notifications/unread_count')
    .then(response => Number(response.data.count ?? 0));
};
