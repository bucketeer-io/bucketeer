import axiosClient from '@api/axios-client';
import { NotificationCreatorResponse } from './notification-creator';
import { NotificationWire, toFeedItem } from './notification-mapper';

export interface NotificationPublisherPayload {
  id: string;
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

// POST /v1/admin_notification/publish — publish an existing draft. System admin only.
export const notificationPublisher = async (
  payload: NotificationPublisherPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .post<NotificationResponseWire>('/v1/admin_notification/publish', payload)
    .then(response => ({
      notification: toFeedItem(response.data.notification)
    }));
};
