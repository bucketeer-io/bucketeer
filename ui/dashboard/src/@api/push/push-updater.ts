import axiosClient from '@api/axios-client';
import { Push } from '@types';

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

export interface PushUpdaterResponse {
  push: Push;
}

export const pushUpdater = async (
  params?: PushUpdaterPayload
): Promise<PushUpdaterResponse> => {
  return axiosClient
    .patch<PushUpdaterResponse>('/v1/push', params)
    .then(response => response.data);
};
