import axiosClient from '@api/axios-client';
import { AuthTypeMap, AuthResponse } from '@types';

export interface ExchangeDemoTokenPayload {
  code: string;
  redirectUrl: string;
  type: AuthTypeMap[keyof AuthTypeMap];
}

export const exchangeDemoToken = async (
  payload: ExchangeDemoTokenPayload
): Promise<AuthResponse> => {
  return axiosClient
    .post<AuthResponse>('/v1/exchange_demo_token', payload)
    .then(response => response.data);
};
