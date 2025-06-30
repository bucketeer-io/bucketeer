import axiosClient from '@api/axios-client';
import { PushResponse } from './push-fetcher';

export interface PushCreatorPayload {
  tags?: string[];
  name: string;
  fcmServiceAccount: Uint8Array | string;
  environmentId: string;
}

export const pushCreator = async (
  payload?: PushCreatorPayload
): Promise<PushResponse> => {
  return axiosClient
    .post<PushResponse>('/v1/push', payload)
    .then(response => response.data);
};
