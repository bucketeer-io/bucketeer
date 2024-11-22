import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface APIKeyDisableParams {
  id: string;
  environmentId: string;
  command: AnyObject;
}

export const apiKeyDisable = async (params?: APIKeyDisableParams) => {
  return axiosClient
    .post('/v1/account/disable_api_key', params)
    .then(response => response.data);
};
