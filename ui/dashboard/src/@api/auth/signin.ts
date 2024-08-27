import axiosClient from '@api/axios-client';
import { AuthResponse, SignInForm } from '@types';

export const signIn = async (formValues: SignInForm): Promise<AuthResponse> => {
  return axiosClient
    .post(`/v1/auth/signin`, formValues)
    .then(response => response.data);
};
