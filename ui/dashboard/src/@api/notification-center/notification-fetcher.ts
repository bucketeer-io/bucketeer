import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';
import { NotificationWire, toFeedItem } from './notification-mapper';

export interface NotificationFetcherParams {
  id: string;
  language?: string;
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

// GET /v1/notification — get a single notification by id. Drafts are visible
// to system admins only.
export const notificationFetcher = async (
  params: NotificationFetcherParams
) => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationResponseWire>(`/v1/notification?${requestParams}`)
    .then(response => toFeedItem(response.data.notification));
};
