import axiosClient from '@api/axios-client';

export interface OrganizationUnarchivePayload {
  id: string;
}

export const organizationUnarchive = async (
  params?: OrganizationUnarchivePayload
) => {
  return axiosClient
    .post('/v1/environment/unarchive_organization', params)
    .then(response => response.data);
};
