import axiosClient from '@api/axios-client';
import { TriggerItemType } from '@types';

export interface TriggerUpdateParams {
  id: string;
  environmentId: string;
  description?: string;
  reset?: boolean;
  disabled?: boolean;
}

export const triggerUpdate = async (
  params?: TriggerUpdateParams
): Promise<TriggerItemType> => {
  return axiosClient
    .patch<TriggerItemType>('/v1/flag_trigger', params)
    .then(response => response.data);
};
