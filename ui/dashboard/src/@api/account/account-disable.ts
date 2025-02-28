import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface AccountDisablePayload {
  email: string;
  organizationId: string;
  command: AnyObject;
}

export const accountDisable = async (params?: AccountDisablePayload) => {
  return axiosClient
    .post('/v1/account/disable_account', params)
    .then(response => response.data);
};
