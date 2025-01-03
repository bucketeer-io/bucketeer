import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface AccountEnableParams {
  email: string;
  organizationId: string;
  command: AnyObject;
}

export const accountEnable = async (params?: AccountEnableParams) => {
  return axiosClient
    .post('/v1/account/enable_account', params)
    .then(response => response.data);
};
