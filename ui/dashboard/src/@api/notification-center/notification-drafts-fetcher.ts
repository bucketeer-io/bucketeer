import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, NotificationCenterDraftCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { NotificationWire, toFeedItem } from 'utils/notification-mapper';
import { stringifyParams } from 'utils/search-params';

export type NotificationDraftsFetcherParams = CollectionParams;

interface DraftCollectionWire {
  notifications: NotificationWire[];
  nextCursor: string;
  totalCount: string;
}

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
