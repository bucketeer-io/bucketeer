import axiosClient from '@api/axios-client';
import {
  NotificationCenterFeedItem,
  NotificationCenterLocalization
} from '@types';
import { NotificationWire, toFeedItem } from 'utils/notification-mapper';

export interface NotificationCreatorPayload {
  localizations: NotificationCenterLocalization[];
}

export interface NotificationCreatorResponse {
  notification: NotificationCenterFeedItem;
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

// POST /v1/admin_notification — create a draft notification. System admin only.
export const notificationCreator = async (
  payload: NotificationCreatorPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .post<NotificationResponseWire>('/v1/admin_notification', payload)
    .then(response => ({
      notification: toFeedItem(response.data.notification)
    }));
};
