import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface APIKeyEnableParams {
  id: string;
  environmentId: string;
  command: AnyObject;
}

export const apiKeyEnable = async (params?: APIKeyEnableParams) => {
  return axiosClient
    .post('/v1/account/enable_api_key', params)
    .then(response => response.data);
};
