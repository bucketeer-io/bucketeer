import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface AccountDisableParams {
  email: string;
  organizationId: string;
  command: AnyObject;
}

export const accountDisable = async (params?: AccountDisableParams) => {
  return axiosClient
    .post('/v1/account/disable_account', params)
    .then(response => response.data);
};
