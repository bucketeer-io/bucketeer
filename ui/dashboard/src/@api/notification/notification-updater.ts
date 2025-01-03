import axiosClient from '@api/axios-client';
import { Notification, SourceType } from '@types';

export interface NotificationUpdterPayload {
  id: string;
  environmentId: string;
  sourceTypes: SourceType[];
  name: string;
  disabled?: boolean;
}

export interface NotificationUpdterResponse {
  account: Array<Notification>;
}

export const notificationUpdater = async (
  payload: NotificationUpdterPayload
): Promise<NotificationUpdterResponse> => {
  return axiosClient
    .patch<NotificationUpdterResponse>('/v1/subscription', payload)
    .then(response => response.data);
};
