import axiosClient from '@api/axios-client';
import { Push } from '@types';

export interface PushUpdaterPayload {
  id: string;
  environmentId?: string;
  description?: string;
  disabled?: boolean;
  name?: string;
  tags?: string[];
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
