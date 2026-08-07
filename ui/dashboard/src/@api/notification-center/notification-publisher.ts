import axiosClient from '@api/axios-client';
import { NotificationWire, toFeedItem } from 'utils/notification-mapper';
import { NotificationCreatorResponse } from './notification-creator';

export interface NotificationPublisherPayload {
  id: string;
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

export const notificationPublisher = async (
  payload: NotificationPublisherPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .post<NotificationResponseWire>('/v1/admin_notification/publish', payload)
    .then(response => ({
      notification: toFeedItem(response.data.notification)
    }));
};
