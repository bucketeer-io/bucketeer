import axiosClient from '@api/axios-client';

export interface OrganizationUpdateParams {
  id: string;
  renameCommand: {
    name: string;
  };
  changeDescriptionCommand: {
    description: string;
  };
}

export const organizationUpdate = async (params?: OrganizationUpdateParams) => {
  return axiosClient
    .post('/v1/environment/update_organization', params)
    .then(response => response.data);
};
