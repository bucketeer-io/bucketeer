import axiosClient from '@api/axios-client';

export interface OrganizationArchiveParams {
  id: string;
}

export const organizationArchive = async (
  params?: OrganizationArchiveParams
) => {
  return axiosClient
    .post('/v1/environment/archive_organization', params)
    .then(response => response.data);
};
