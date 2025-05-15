import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { TriggerItemType } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface TriggerDeleteParams {
  environmentId: string;
  id: string;
}

export const triggerDelete = async (
  _params?: TriggerDeleteParams
): Promise<TriggerItemType> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .delete<TriggerItemType>(`/v1/flag_trigger?${stringifyParams(params)}`)
    .then(response => response.data);
};
