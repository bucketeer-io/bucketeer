import axiosClient from '@api/axios-client';
import { Push } from '@types';

type Tags = {
  tags?: string[];
};

type Name = {
  name?: string;
};

export type PushUpdaterParams = Tags &
  Name & {
    id: string;
    environmentId?: string;
    description?: string;
    disabled?: boolean;
    addPushTagsCommand?: Tags;
    deletePushTagsCommand?: Tags;
    renamePushCommand?: Name;
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
