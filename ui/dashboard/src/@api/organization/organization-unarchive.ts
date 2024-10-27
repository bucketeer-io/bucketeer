import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface OrganizationUnArchiveParams {
  id: string;
  command: AnyObject;
}

export const organizationUnArchive = async (
  params?: OrganizationUnArchiveParams
) => {
  return axiosClient
    .post('/v1/environment/unarchive_organization', params)
    .then(response => response.data);
};
