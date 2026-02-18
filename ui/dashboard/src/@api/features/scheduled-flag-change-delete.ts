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
) => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .delete(`/v1/scheduled_flag_change?${requestParams}`)
    .then(response => response.data);
};
