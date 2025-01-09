import axiosClient from '@api/axios-client';

export interface AccountDisableParams {
  email: string;
  organizationId: string;
}

export const accountDisable = async (params?: AccountDisableParams) => {
  return axiosClient
    .post('/v1/account/disable_account', params)
    .then(response => response.data);
};
