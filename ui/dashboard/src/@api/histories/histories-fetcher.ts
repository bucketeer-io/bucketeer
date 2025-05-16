import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AuditLogCollection, CollectionParams } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface HistoriesFetcherParams extends CollectionParams {
  from?: string;
  to?: string;
  entityType?: number;
  environmentId?: string;
  featureId: string;
}

export const historiesFetcher = async (
  params?: HistoriesFetcherParams
): Promise<AuditLogCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<AuditLogCollection>(`/v1/feature_history?${requestParams}`)
    .then(response => response.data);
};
