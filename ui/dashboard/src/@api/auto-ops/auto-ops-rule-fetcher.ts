import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { AutoOpsRule } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AutoOpsRuleFetcherParams {
  id: string;
  environmentId: string;
}

export interface AutoOpsRuleFetcherResponse {
  autoOpsRule: AutoOpsRule;
}

export const autoOpsRuleFetcher = async (
  _params?: AutoOpsRuleFetcherParams
): Promise<AutoOpsRuleFetcherResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<AutoOpsRuleFetcherResponse>(
      `/v1/auto_ops_rule?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
