import axiosClient from '@api/axios-client';

export interface AccountEnableParams {
  email: string;
  organizationId: string;
}

export const accountEnable = async (params?: AccountEnableParams) => {
  return axiosClient
    .post('/v1/account/enable_account', params)
    .then(response => response.data);
};
