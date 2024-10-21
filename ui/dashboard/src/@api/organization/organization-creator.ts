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

export interface OrganizationCreatorParams {
  command: OrganizationCreatorCommand;
}

export interface OrganizationResponse {
  organization: Array<Organization>;
}

export const organizationCreator = async (
  params?: OrganizationCreatorParams
): Promise<OrganizationResponse> => {
  return axiosClient
    .post<OrganizationResponse>('/v1/environment/create_organization', params)
    .then(response => response.data);
};
