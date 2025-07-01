import axiosClient from '@api/axios-client';
import { NotificationRecipient, SourceType } from '@types';
import { NotificationResponse } from './notification-fetcher';

export interface NotificationCreatorPayload {
  environmentId: string;
  name: string;
  sourceTypes: SourceType[];
  recipient: NotificationRecipient;
  featureFlagTags: string[];
}

export const notificationCreator = async (
  payload: NotificationCreatorPayload
): Promise<NotificationResponse> => {
  return axiosClient
    .post<NotificationResponse>('/v1/subscription', payload)
    .then(response => response.data);
};
