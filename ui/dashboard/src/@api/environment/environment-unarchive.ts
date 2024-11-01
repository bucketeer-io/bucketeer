import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface EnvironmentUnArchiveParams {
  id: string;
  command: AnyObject;
}

export const environmentUnArchive = async (
  params?: EnvironmentUnArchiveParams
) => {
  return axiosClient
    .post('/v1/environment/unarchive_environment', params)
    .then(response => response.data);
};
