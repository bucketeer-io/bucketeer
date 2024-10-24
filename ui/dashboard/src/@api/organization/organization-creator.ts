import axiosClient from '@api/axios-client';
import { Organization } from '@types';

export interface OrganizationCreatorCommand {
  name: string;
  urlCode: string;
  description?: string;
  isTrial?: boolean;
  isSystemAdmin: boolean;
  ownerEmail: string;
}

export interface OrganizationResponse {
  organization: Array<Organization>;
}

export const organizationCreator = async (
  params?: OrganizationCreatorCommand
): Promise<OrganizationResponse> => {
  return axiosClient
    .post<OrganizationResponse>('/v1/environment/create_organization', params)
    .then(response => response.data);
};
