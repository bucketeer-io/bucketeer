import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, NotificationCenterFeedCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationsFetcherParams extends CollectionParams {
  environmentId?: string;
  read?: boolean;
  from?: string;
  to?: string;
}

// GET /v1/notifications — list published notifications.
export const notificationsFetcher = async (
  params?: NotificationsFetcherParams
): Promise<NotificationCenterFeedCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationCenterFeedCollection>(`/v1/notifications?${requestParams}`)
    .then(response => response.data);
};
