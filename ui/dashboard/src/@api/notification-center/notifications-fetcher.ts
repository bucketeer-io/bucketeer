import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, NotificationCenterFeedCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';
import { NotificationWire, toFeedItem } from './notification-mapper';

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

// GET /v1/notifications — list published notifications.
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
