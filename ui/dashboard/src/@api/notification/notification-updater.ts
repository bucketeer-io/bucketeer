import axiosClient from '@api/axios-client';
import { Notification, SourceType } from '@types';

export interface NotificationCreatorPayload {
  id: string;
  environmentId: string;
  sourceTypes: SourceType[];
  name: string;
  disabled?: boolean;
}

export interface NotificationCreatorResponse {
  account: Array<Notification>;
}

export const notificationUpdater = async (
  payload: NotificationCreatorPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .patch<NotificationCreatorResponse>('/v1/subscription', payload)
    .then(response => response.data);
};
