import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AutoOpsGoalUserCountFetcherParams {
  environmentId: string;
  featureId: string;
  featureVersion: number;
  opsRuleId: string;
  clauseId: string;
  variationId: string;
}

export interface AutoOpsGoalUserCountResponse {
  opsRuleId: string;
  clauseId: string;
  count: string;
}

export const autoOpsGoalUserCountFetcher = async (
  _params?: AutoOpsGoalUserCountFetcherParams
): Promise<AutoOpsGoalUserCountResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<AutoOpsGoalUserCountResponse>(
      `/v1/ops_goal_user_count?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
