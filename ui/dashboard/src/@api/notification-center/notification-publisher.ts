import axiosClient from '@api/axios-client';
import { NotificationCenterLocalization } from '@types';
import { NotificationCreatorResponse } from './notification-creator';

export interface NotificationPublisherPayload {
  environmentId: string;
  // Set when publishing an existing draft in place; omitted to publish a
  // brand-new notification directly.
  id?: string;
  localizations: NotificationCenterLocalization[];
}

// POST /v1/notification/publish — publish a draft. System admin only.
export const notificationPublisher = async (
  payload: NotificationPublisherPayload
): Promise<NotificationCreatorResponse> => {
  return axiosClient
    .post<NotificationCreatorResponse>('/v1/notification/publish', payload)
    .then(response => response.data);
};
