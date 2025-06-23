import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { Notification } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationFetcherPayload {
  environmentId: string;
  id: string;
}

export interface NotificationFetcherResponse {
  subscription: Notification;
}

export const notificationFetcher = async (
  payload?: NotificationFetcherPayload
): Promise<NotificationFetcherResponse> => {
  const params = stringifyParams(pickBy(payload, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationFetcherResponse>(`/v1/subscription?${params}`)
    .then(response => response.data);
};
