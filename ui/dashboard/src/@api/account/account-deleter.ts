import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface AccountDeleterParams {
  email: string;
  organizationId: string;
  command: AnyObject;
}

export const accountDeleter = async (params?: AccountDeleterParams) => {
  return axiosClient
    .post('/v1/account/delete_account', params)
    .then(response => response.data);
};
