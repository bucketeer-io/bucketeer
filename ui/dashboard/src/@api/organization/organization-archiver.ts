import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface OrganizationArchiverParams {
  id: string;
  command: AnyObject;
}

export const organizationArchiver = async (
  params?: OrganizationArchiverParams
) => {
  return axiosClient
    .post('/v1/environment/archive_organization', params)
    .then(response => response.data);
};
