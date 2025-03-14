import axiosClient from '@api/axios-client';
import {
  AutoOpsRule,
  AutoOpsType,
  DatetimeClause,
  OpsEventRateClause
} from '@types';

export interface AutoOpsCreatorParams {
  environmentId: string;
  featureId: string;
  opsType: AutoOpsType;
  opsEventRateClauses?: OpsEventRateClause[];
  datetimeClauses?: DatetimeClause[];
}

export interface AutoOpsCreatorResponse {
  autoOpsRule: AutoOpsRule;
}

export const autoOpsCreator = async (
  params?: AutoOpsCreatorParams
): Promise<AutoOpsCreatorResponse> => {
  return axiosClient
    .post<AutoOpsCreatorResponse>('/v1/auto_ops_rule', params)
    .then(response => response.data);
};
