import axiosClient from '@api/axios-client';
import { Notification, NotificationRecipient, SourceType } from '@types';

export interface NotificationCreatorPayload {
  environmentId: string;
  name: string;
  sourceTypes: SourceType[];
  recipient: NotificationRecipient;
}

export interface NotificationCreatorResponse {
  subscription: Notification;
}

export const notificationCreator = async (
  payload: NotificationCreatorPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .post<NotificationCreatorResponse>('/v1/subscription', payload)
    .then(response => response.data);
};
