import axiosClient from '@api/axios-client';
import { DatetimeClause, OpsEventRateClause } from '@types';
import { AutoOpsCreatorResponse } from './auto-ops-creator';

export interface ClauseUpdateType<T> {
  id: string;
  delete: boolean;
  clause: T;
}

export interface AutoOpsUpdateParams {
  id: string;
  environmentId: string;
  updateOpsEventRateClauses?: ClauseUpdateType<OpsEventRateClause>[];
  updateDatetimeClauses?: ClauseUpdateType<DatetimeClause>[];
}

export const autoOpsUpdate = async (
  params?: AutoOpsUpdateParams
): Promise<AutoOpsCreatorResponse> => {
  return axiosClient
    .patch<AutoOpsCreatorResponse>('/v1/auto_ops_rule', params)
    .then(response => response.data);
};
