import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, NotificationCenterFeedCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { NotificationWire, toFeedItem } from 'utils/notification-mapper';
import { stringifyParams } from 'utils/search-params';

export type NotificationReadStatus = 'ALL' | 'UNREAD' | 'READ';

export interface NotificationsFetcherParams extends CollectionParams {
  readStatus?: NotificationReadStatus;
  publishedAtFrom?: string;
  publishedAtTo?: string;
  language?: string;
}

interface FeedCollectionWire {
  notifications: NotificationWire[];
  nextCursor: string;
  totalCount: string;
}

export const notificationsFetcher = async (
  params?: NotificationsFetcherParams
): Promise<NotificationCenterFeedCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<FeedCollectionWire>(`/v1/notifications?${requestParams}`)
    .then(response => ({
      notifications: response.data.notifications.map(toFeedItem),
      cursor: response.data.nextCursor,
      totalCount: response.data.totalCount
    }));
};
