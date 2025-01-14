import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { Account } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AccountByOrgFetcherParams {
  email: string;
  organizationId: string;
}

export interface AccountByOrgFetcherResponse {
  account: Account;
}

export const accountByOrgFetcher = async (
  _params?: AccountByOrgFetcherParams
): Promise<AccountByOrgFetcherResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<AccountByOrgFetcherResponse>(
      `/v1/account/get_account?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
