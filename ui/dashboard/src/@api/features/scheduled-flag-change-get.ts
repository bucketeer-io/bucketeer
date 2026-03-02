import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import type { GetScheduledFlagChangeResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ScheduledFlagChangeGetParams {
  environmentId: string;
  id: string;
}

export const scheduledFlagChangeGet = async (
  params?: ScheduledFlagChangeGetParams
): Promise<GetScheduledFlagChangeResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<GetScheduledFlagChangeResponse>(
      `/v1/scheduled_flag_change?${requestParams}`
    )
    .then(response => response.data);
};
