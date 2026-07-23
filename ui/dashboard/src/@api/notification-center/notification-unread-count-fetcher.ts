import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { NotificationCenterUnreadCount } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationUnreadCountFetcherParams {
  environmentId?: string;
}

// GET /v1/notifications/unread_count — unread count for the bell badge.
export const notificationUnreadCountFetcher = async (
  params?: NotificationUnreadCountFetcherParams
): Promise<NotificationCenterUnreadCount> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationCenterUnreadCount>(
      `/v1/notifications/unread_count?${requestParams}`
    )
    .then(response => response.data);
};
