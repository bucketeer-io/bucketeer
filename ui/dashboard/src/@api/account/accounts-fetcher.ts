import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AccountCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AccountsFetcherParams extends CollectionParams {
  organizationId?: string;
  organizationRole?: number;
  environmentId?: string;
  environmentRole?: number;
}

export const accountsFetcher = async (
  _params?: AccountsFetcherParams
): Promise<AccountCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<AccountCollection>(
      `/v1/account/list_accounts?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
