import axiosClient from '@api/axios-client';
import { ConsoleAccountResponse } from '@types';

export interface MeFetcherPayload {
  organizationId: string;
}

export const accountMeFetcher = async (
  payload: MeFetcherPayload
): Promise<ConsoleAccountResponse> => {
  return axiosClient
    .post<ConsoleAccountResponse>('/v1/account/get_me', payload)
    .then(response => response.data);
};
