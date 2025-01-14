import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { ConsoleAccountResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AccountByEnvFetcherParams {
  email: string;
  environmentId: string;
}

export const accountByEnvFetcher = async (
  _params: AccountByEnvFetcherParams
): Promise<ConsoleAccountResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<ConsoleAccountResponse>(
      `/v1/account/get_account_by_environment?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
