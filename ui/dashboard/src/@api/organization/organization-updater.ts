import axiosClient from '@api/axios-client';

export interface OrganizationUpdateParams {
  id: string;
  name: string;
  description?: string;
  ownerEmail: string;
}

export const organizationUpdater = async (
  params?: OrganizationUpdateParams
) => {
  return axiosClient
    .post('/v1/environment/update_organization', params)
    .then(response => response.data);
};
