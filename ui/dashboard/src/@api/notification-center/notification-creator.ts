import axiosClient from '@api/axios-client';
import {
  NotificationCenterFeedItem,
  NotificationCenterPublishPayload
} from '@types';

export interface NotificationCreatorPayload extends NotificationCenterPublishPayload {
  environmentId: string;
}

export interface NotificationCreatorResponse {
  notification: NotificationCenterFeedItem;
}

// POST /v1/notification — create a draft notification. System admin only.
export const notificationCreator = async (
  payload: NotificationCreatorPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .post<NotificationCreatorResponse>('/v1/notification', payload)
    .then(response => response.data);
};
