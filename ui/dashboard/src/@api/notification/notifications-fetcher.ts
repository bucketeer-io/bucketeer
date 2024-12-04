import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AccountCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface NotificationsFetcherParams extends CollectionParams {
  environmentId?: string;
}

export const notificationsFetcher = async (
  params?: NotificationsFetcherParams
): Promise<AccountCollection> => {
  return axiosClient
    .post<AccountCollection>(
      '/v1/account/list_accounts',
      pickBy(params, v => isNotEmpty(v))
    )
    .then(response => response.data);
};
