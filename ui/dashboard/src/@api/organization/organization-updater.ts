import axiosClient from '@api/axios-client';

export interface OrganizationUpdateParams {
  id: string;
  renameCommand: {
    name: string;
  };
  changeDescriptionCommand: {
    description?: string;
  };
  changeOwnerEmailCommand: {
    ownerEmail: string;
  };
}

export const organizationUpdater = async (
  params?: OrganizationUpdateParams
) => {
  return axiosClient
    .post('/v1/environment/update_organization', params)
    .then(response => response.data);
};
