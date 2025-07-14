import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ScheduleDeleteParams {
  id: string;
  environmentId: string;
}

export const scheduleFlagDelete = async (params?: ScheduleDeleteParams) => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .delete(`/v1/schedule_flag?${requestParams}`)
    .then(response => response.data);
};
