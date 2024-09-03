import axiosClient from '@api/axios-client';
import { AuthTypeMap, AuthResponse } from '@types';

export interface ExchangeTokenPayload {
  code: string;
  redirectUrl: string;
  type: AuthTypeMap[keyof AuthTypeMap];
}

export const exchangeToken = async (
  payload: ExchangeTokenPayload
): Promise<AuthResponse> => {
  return axiosClient
    .post<AuthResponse>('/v1/auth/exchange_token', payload)
    .then(response => response.data);
};
