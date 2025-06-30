import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface OrganizationArchiveParams {
  id: string;
  command?: AnyObject;
}

export const organizationArchive = async (
  params?: OrganizationArchiveParams
) => {
  return axiosClient
    .post('/v1/environment/archive_organization', params)
    .then(response => response.data);
};
