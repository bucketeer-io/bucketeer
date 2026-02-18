import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import type {
  ScheduledFlagChangeCollection,
  ScheduledFlagChangeStatus
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export type ScheduledFlagChangesOrderBy =
  | 'DEFAULT'
  | 'SCHEDULED_AT'
  | 'CREATED_AT';
export type ScheduledFlagChangesOrderDirection = 'ASC' | 'DESC';

export interface ScheduledFlagChangesFetcherParams {
  environmentId: string;
  featureId?: string;
  statuses?: ScheduledFlagChangeStatus[];
  fromScheduledAt?: string;
  toScheduledAt?: string;
  pageSize?: number;
  cursor?: string;
  orderBy?: ScheduledFlagChangesOrderBy;
  orderDirection?: ScheduledFlagChangesOrderDirection;
}

export const scheduledFlagChangesFetcher = async (
  params?: ScheduledFlagChangesFetcherParams
): Promise<ScheduledFlagChangeCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<ScheduledFlagChangeCollection>(
      `/v1/scheduled_flag_changes?${requestParams}`
    )
    .then(response => response.data);
};
