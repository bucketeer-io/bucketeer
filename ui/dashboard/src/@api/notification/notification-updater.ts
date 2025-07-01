import axiosClient from '@api/axios-client';
import { NotificationLanguage, SourceType } from '@types';
import { NotificationResponse } from './notification-fetcher';

export interface NotificationUpdterPayload {
  id: string;
  environmentId?: string;
  sourceTypes?: SourceType[];
  name?: string;
  disabled?: boolean;
  language?: NotificationLanguage;
  featureFlagTags?: string[];
}

export const notificationUpdater = async (
  payload: NotificationUpdterPayload
): Promise<NotificationResponse> => {
  return axiosClient
    .patch<NotificationResponse>('/v1/subscription', payload)
    .then(response => response.data);
};
