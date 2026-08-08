import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, NotificationCenterDraftCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';
import { NotificationWire, toFeedItem } from './notification-mapper';

export type NotificationDraftsFetcherParams = CollectionParams;

interface DraftCollectionWire {
  notifications: NotificationWire[];
  nextCursor: string;
  totalCount: string;
}

// GET /v1/admin_notifications/drafts — list draft notifications. System admin only.
export const notificationDraftsFetcher = async (
  params?: NotificationDraftsFetcherParams
): Promise<NotificationCenterDraftCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<DraftCollectionWire>(`/v1/admin_notifications/drafts?${requestParams}`)
    .then(response => ({
      notifications: response.data.notifications.map(toFeedItem),
      cursor: response.data.nextCursor,
      totalCount: response.data.totalCount
    }));
};
