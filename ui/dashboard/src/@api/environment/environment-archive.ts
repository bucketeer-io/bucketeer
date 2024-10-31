import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface EnvironmentArchiveParams {
  id: string;
  command: AnyObject;
}

export const environmentArchive = async (params?: EnvironmentArchiveParams) => {
  return axiosClient
    .post('/v1/environment/archive_environment', params)
    .then(response => response.data);
};
