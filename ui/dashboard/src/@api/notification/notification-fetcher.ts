import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { Notification } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface NotificationFetcherPayload {
  id: string;
  environmentId: string;
}

export interface NotificationFetcherResponse {
  subscription: Notification;
}

export const notificationFetcher = async (
  payload: NotificationFetcherPayload
): Promise<NotificationFetcherResponse> => {
  const params = pickBy(payload, v => isNotEmpty(v));

  return axiosClient
    .get<NotificationFetcherResponse>('/v1/subscription', {
      params
    })
    .then(response => response.data);
};
