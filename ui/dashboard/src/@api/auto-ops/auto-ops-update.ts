import axiosClient from '@api/axios-client';
import { DatetimeClause, OpsEventRateClause } from '@types';
import { AutoOpsCreatorResponse } from './auto-ops-creator';

interface ClauseUpdateType<T> {
  id: string;
  delete: boolean;
  clause: T[];
}

export interface AutoOpsUpdateParams {
  environmentId: string;
  updateOpsEventRateClauses?: ClauseUpdateType<OpsEventRateClause>;
  updateDatetimeClauses?: ClauseUpdateType<DatetimeClause>;
}

export const autoOpsUpdate = async (
  params?: AutoOpsUpdateParams
): Promise<AutoOpsCreatorResponse> => {
  return axiosClient
    .patch<AutoOpsCreatorResponse>('/v1/auto_ops_rule', params)
    .then(response => response.data);
};
