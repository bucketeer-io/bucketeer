import axiosClient from '@api/axios-client';
import { AccountCollection, CollectionParams } from '@types';

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
    .post<AccountCollection>('/v1/account/list_accounts_v2', params)
    .then(response => response.data);
};
