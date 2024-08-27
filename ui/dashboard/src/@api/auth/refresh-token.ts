import axiosClient from '@api/axios-client';
import { AuthResponse } from '@types';

export const refreshTokenFetcher = async (
  token: string
): Promise<AuthResponse> => {
  return axiosClient
    .post(`/v1/auth/refresh_token`, { refreshToken: token })
    .then(response => response.data);
};
