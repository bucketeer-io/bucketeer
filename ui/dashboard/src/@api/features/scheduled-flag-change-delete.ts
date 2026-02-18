import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ScheduledFlagChangeDeleteParams {
  environmentId: string;
  id: string;
}

export const scheduledFlagChangeDelete = async (
  params?: ScheduledFlagChangeDeleteParams
): Promise<void> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .delete<void>(`/v1/scheduled_flag_change?${requestParams}`)
    .then(() => {});
};
