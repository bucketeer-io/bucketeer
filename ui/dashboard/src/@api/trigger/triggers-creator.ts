import axiosClient from '@api/axios-client';
import { TriggerActionType, TriggerItemType, TriggerType } from '@types';

export interface TriggerCreatorParams {
  environmentId: string;
  featureId: string;
  type: TriggerType;
  action: TriggerActionType;
  description?: string;
}

export const triggerCreator = async (
  params?: TriggerCreatorParams
): Promise<TriggerItemType> => {
  return axiosClient
    .post<TriggerItemType>('/v1/flag_trigger', params)
    .then(response => response.data);
};
