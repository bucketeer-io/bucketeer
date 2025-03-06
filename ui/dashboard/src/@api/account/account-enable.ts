import axiosClient from '@api/axios-client';

export interface AccountEnablePayload {
  email: string;
  organizationId: string;
}

export const accountEnable = async (params?: AccountEnablePayload) => {
  return axiosClient
    .post('/v1/account/enable_account', params)
    .then(response => response.data);
};
