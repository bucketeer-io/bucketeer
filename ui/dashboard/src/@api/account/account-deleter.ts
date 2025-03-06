import axiosClient from '@api/axios-client';

export interface AccountDeleterPayload {
  email: string;
  organizationId: string;
}

export const accountDeleter = async (params?: AccountDeleterPayload) => {
  return axiosClient
    .post('/v1/account/delete_account', params)
    .then(response => response.data);
};
