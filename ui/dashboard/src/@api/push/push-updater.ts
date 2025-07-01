import axiosClient from '@api/axios-client';
import { PushResponse } from './push-fetcher';

export type TagChangeActionType =
  | 'UNSPECIFIED'
  | 'CREATE'
  | 'UPDATE'
  | 'DELETE';
export interface TagChange {
  changeType: TagChangeActionType;
  tag: string;
}
export interface PushUpdaterPayload {
  id: string;
  environmentId?: string;
  description?: string;
  disabled?: boolean;
  name?: string;
  tags?: string[];
  tagChanges?: TagChange[];
}

export const pushUpdater = async (
  params?: PushUpdaterPayload
): Promise<PushResponse> => {
  return axiosClient
    .patch<PushResponse>('/v1/push', params)
    .then(response => response.data);
};
