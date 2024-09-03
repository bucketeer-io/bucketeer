import axiosClient from '@api/axios-client';
import { AuthTypeMap, AuthUrlResponse } from '@types';

export interface AuthenticationUrlPayload {
  state: string;
  redirectUrl: string;
  type: AuthTypeMap[keyof AuthTypeMap];
}

export const authenticationUrl = async (
  payload: AuthenticationUrlPayload
): Promise<AuthUrlResponse> => {
  return axiosClient
    .post<AuthUrlResponse>('/v1/auth/authentication_url', payload)
    .then(response => response.data);
};
