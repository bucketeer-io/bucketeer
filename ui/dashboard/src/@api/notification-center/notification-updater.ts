import axiosClient from '@api/axios-client';
import { NotificationCenterUpdatePayload } from '@types';
import { NotificationCreatorResponse } from './notification-creator';

export interface NotificationUpdaterPayload extends NotificationCenterUpdatePayload {
  environmentId: string;
}

// PATCH /v1/notification — update a draft's localizations. System admin only.
export const notificationUpdater = async (
  payload: NotificationUpdaterPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .patch<NotificationCreatorResponse>('/v1/notification', payload)
    .then(response => response.data);
};
