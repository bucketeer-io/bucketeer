import axiosClient from '@api/axios-client';
import { NotificationCenterLocalization } from '@types';
import { NotificationWire, toFeedItem } from 'utils/notification-mapper';
import { NotificationCreatorResponse } from './notification-creator';

export interface NotificationUpdaterPayload {
  id: string;
  localizations: NotificationCenterLocalization[];
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

export const notificationUpdater = async (
  payload: NotificationUpdaterPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .patch<NotificationResponseWire>('/v1/admin_notification', payload)
    .then(response => ({
      notification: toFeedItem(response.data.notification)
    }));
};
