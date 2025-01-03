import axiosClient from '@api/axios-client';
import { Push } from '@types';

export type PushCreatorParams = {
  tags: string[];
  name: string;
  fcmServiceAccount: string;
  environmentId: string;
};

export interface PushCreatorResponse {
  apiKey: Array<Push>;
}

export const pushCreator = async (
  params?: PushCreatorParams
): Promise<PushCreatorResponse> => {
  return axiosClient
    .post<PushCreatorResponse>('/v1/push', params)
    .then(response => response.data);
};
