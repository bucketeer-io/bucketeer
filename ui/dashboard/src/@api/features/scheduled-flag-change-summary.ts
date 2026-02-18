import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import type { GetScheduledFlagChangeSummaryResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ScheduledFlagChangeSummaryParams {
  environmentId: string;
  featureId: string;
}

export const scheduledFlagChangeSummaryFetcher = async (
  params?: ScheduledFlagChangeSummaryParams
): Promise<GetScheduledFlagChangeSummaryResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<GetScheduledFlagChangeSummaryResponse>(
      `/v1/scheduled_flag_change/summary?${requestParams}`
    )
    .then(response => response.data);
};
