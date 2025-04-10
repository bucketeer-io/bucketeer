import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AuditLogCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AuditLogsFetcherParams extends CollectionParams {
  from?: string;
  to?: string;
  entityType?: number;
  environmentId?: string;
}

export const auditLogsFetcher = async (
  params?: AuditLogsFetcherParams
): Promise<AuditLogCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<AuditLogCollection>(`/v1/audit_logs?${requestParams}`)
    .then(response => response.data);
};
