import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, AutoOpsCountCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AutoOpsCountFetcherParams extends CollectionParams {
  environmentId: string;
  featureIds?: string[];
  autoOpsRuleIds?: string[];
}

export const autoOpsCountFetcher = async (
  _params?: AutoOpsCountFetcherParams
): Promise<AutoOpsCountCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<AutoOpsCountCollection>(
      `/v1/auto_ops_rule/ops_counts?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
