import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { NotificationWire, toFeedItem } from 'utils/notification-mapper';
import { stringifyParams } from 'utils/search-params';

export interface NotificationFetcherParams {
  id: string;
  language?: string;
}

interface NotificationResponseWire {
  notification: NotificationWire;
}

export const notificationFetcher = async (
  params: NotificationFetcherParams
) => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationResponseWire>(`/v1/notification?${requestParams}`)
    .then(response => toFeedItem(response.data.notification));
};
