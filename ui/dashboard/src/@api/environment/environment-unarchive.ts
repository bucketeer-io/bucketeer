import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface EnvironmentUnarchiveParams {
  id: string;
  command: AnyObject;
}

export const environmentUnarchive = async (
  params?: EnvironmentUnarchiveParams
) => {
  return axiosClient
    .post('/v1/environment/unarchive_environment', params)
    .then(response => response.data);
};
