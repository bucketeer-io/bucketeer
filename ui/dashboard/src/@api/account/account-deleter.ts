import axiosClient from '@api/axios-client';

export interface AccountDeleterParams {
  email: string;
  organizationId: string;
}

export const accountDeleter = async (params?: AccountDeleterParams) => {
  return axiosClient
    .post('/v1/account/delete_account', params)
    .then(response => response.data);
};
