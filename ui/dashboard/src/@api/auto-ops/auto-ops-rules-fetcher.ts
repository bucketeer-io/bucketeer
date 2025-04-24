import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, AutoOpsRuleCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AutoOpsRulesFetcherParams extends CollectionParams {
  environmentId: string;
  featureIds?: string[];
}

export const autoOpsRulesFetcher = async (
  _params?: AutoOpsRulesFetcherParams
): Promise<AutoOpsRuleCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<AutoOpsRuleCollection>(`/v1/auto_ops_rules?${stringifyParams(params)}`)
    .then(response => response.data);
};
