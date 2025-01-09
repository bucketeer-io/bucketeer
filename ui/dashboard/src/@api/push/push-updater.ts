import axiosClient from '@api/axios-client';
import { Push } from '@types';

export type PushUpdaterParams = {
  id: string;
  environmentId?: string;
  description?: string;
  disabled?: boolean;
  name?: string;
  tags?: string[];
};

export interface PushUpdaterResponse {
  push: Array<Push>;
}

export const pushUpdater = async (
  params?: PushUpdaterParams
): Promise<PushUpdaterResponse> => {
  return axiosClient
    .patch<PushUpdaterResponse>('/v1/push', params)
    .then(response => response.data);
};
