import axiosClient from '@api/axios-client';

export interface OrganizationUnarchiveParams {
  id: string;
}

export const organizationUnarchive = async (
  params?: OrganizationUnarchiveParams
) => {
  return axiosClient
    .post('/v1/environment/unarchive_organization', params)
    .then(response => response.data);
};
