import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { ConsoleAccountResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface MeFetcherParams {
  organizationId: string;
}

export const accountMeFetcher = async (
  _params: MeFetcherParams
): Promise<ConsoleAccountResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<ConsoleAccountResponse>(
      `/v1/account/get_me?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
