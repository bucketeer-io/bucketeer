import axiosClient from '@api/axios-client';
import { Organization } from '@types';

export interface OrganizationDemoCreatorPayload {
  name: string;
  urlCode: string;
  description?: string;
  ownerEmail: string;
}

export interface OrganizationDemoResponse {
  organization: Organization;
}

export const organizationDemoCreator = async (
  params?: OrganizationDemoCreatorPayload
): Promise<OrganizationDemoResponse> => {
  return axiosClient
    .post<OrganizationDemoResponse>(
      '/v1/environment/create_demo_organization',
      params
    )
    .then(response => response.data);
};
