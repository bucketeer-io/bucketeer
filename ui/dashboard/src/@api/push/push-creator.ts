import axiosClient from '@api/axios-client';
import { Push } from '@types';

export type PushCreatorPayload = {
  tags: string[];
  name: string;
  fcmServiceAccount: Uint8Array | string;
  environmentId: string;
};

export interface PushCreatorResponse {
  apiKey: Array<Push>;
}

export const pushCreator = async (
  payload?: PushCreatorPayload
): Promise<PushCreatorResponse> => {
  return axiosClient
    .post<PushCreatorResponse>('/v1/push', payload)
    .then(response => response.data);
};
