import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { NotificationsCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifySearchParams } from 'utils/search-params';

export interface NotificationsFetcherParams extends CollectionParams {
  environmentId?: string;
  organizationId: string;
}

export const notificationsFetcher = async (
  params?: NotificationsFetcherParams
): Promise<NotificationsCollection> => {
  const requestParams = stringifySearchParams(
    pickBy(params, v => isNotEmpty(v))
  );

  return axiosClient
    .get<NotificationsCollection>(`/v1/subscriptions?${requestParams}`)
    .then(response => response.data);
};
