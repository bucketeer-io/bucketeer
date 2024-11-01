import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AccountCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface AccountsFetcherParams extends CollectionParams {
  organizationId?: string;
  organizationRole?: number;
  environmentId?: string;
  environmentRole?: number;
}

export const accountsFetcher = async (
  params?: AccountsFetcherParams
): Promise<AccountCollection> => {
  return axiosClient
    .post<AccountCollection>(
      '/v1/account/list_accounts',
      pickBy(params, v => isNotEmpty(v))
    )
    .then(response => response.data);
};
