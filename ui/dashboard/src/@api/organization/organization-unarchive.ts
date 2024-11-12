import axiosClient from '@api/axios-client';
import { AnyObject } from 'yup';

export interface OrganizationUnarchiveParams {
  id: string;
  command: AnyObject;
}

export const organizationUnarchive = async (
  params?: OrganizationUnarchiveParams
) => {
  return axiosClient
    .post('/v1/environment/unarchive_organization', params)
    .then(response => response.data);
};
