import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AuditLog } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AuditLogDetailsFetcherParams {
  id: string;
  environmentId: string;
}

export interface AuditLogDetailsResponse {
  auditLog: AuditLog;
}

export const auditLogDetailsFetcher = async (
  params?: AuditLogDetailsFetcherParams
): Promise<AuditLogDetailsResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<AuditLogDetailsResponse>(`/v1/audit_log?${requestParams}`)
    .then(response => response.data);
};
