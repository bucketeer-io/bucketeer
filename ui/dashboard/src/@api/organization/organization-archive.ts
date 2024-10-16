import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface ArchiveParams {
  id: string;
  command: AnyObject;
}

export const organizationArchive = async (params?: ArchiveParams) => {
  return axiosClient
    .post('/v1/environment/archive_organization', params)
    .then(response => response.data);
};
