import axiosClient from '@api/axios-client';

export interface AccountDisablePayload {
  email: string;
  organizationId: string;
}

export const accountDisable = async (params?: AccountDisablePayload) => {
  return axiosClient
    .post('/v1/account/disable_account', params)
    .then(response => response.data);
};
