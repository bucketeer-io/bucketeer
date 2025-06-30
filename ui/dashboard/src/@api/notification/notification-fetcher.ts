import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { Notification } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationFetcherPayload {
  environmentId: string;
  id: string;
}

export interface NotificationResponse {
  subscription: Notification;
}

export const notificationFetcher = async (
  payload?: NotificationFetcherPayload
): Promise<NotificationResponse> => {
  const params = stringifyParams(pickBy(payload, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationResponse>(`/v1/subscription?${params}`)
    .then(response => response.data);
};
