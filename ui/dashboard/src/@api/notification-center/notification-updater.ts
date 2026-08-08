import axiosClient from '@api/axios-client';
import { NotificationCenterLocalization } from '@types';
import { NotificationCreatorResponse } from './notification-creator';
import { NotificationWire, toFeedItem } from './notification-mapper';

export interface NotificationUpdaterPayload {
  id: string;
  localizations: NotificationCenterLocalization[];
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

// PATCH /v1/admin_notification — update a draft's localizations. System admin only.
export const notificationUpdater = async (
  payload: NotificationUpdaterPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .patch<NotificationResponseWire>('/v1/admin_notification', payload)
    .then(response => ({
      notification: toFeedItem(response.data.notification)
    }));
};
